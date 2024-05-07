package main

import (
	"bytes"
	"log"
	"log/slog"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"

	"github.com/madkins23/go-slog/internal/data"
)

// -----------------------------------------------------------------------------

type handlerCoords struct {
	x, y float64
}

func (hc *handlerCoords) adjust(by float64) *handlerCoords {
	return &handlerCoords{
		x: hc.x * by,
		y: hc.y * by,
	}
}

type sizeData struct {
	low    handlerCoords
	adjust *handlerCoords
}

var sizes = []*sizeData{
	{
		adjust: defaultLabelSize,
	},
	{
		low: handlerCoords{
			x: 25,
			y: 25,
		},
		adjust: defaultLabelSize.adjust(0.7),
	},
	{
		low: handlerCoords{
			x: 50,
			y: 50,
		},
		adjust: defaultLabelSize.adjust(0.5),
	},
	{
		low: handlerCoords{
			x: 75,
			y: 75,
		},
		adjust: defaultLabelSize.adjust(0.4),
	},
}

// scoreFunction generates an SVG chart for the score visualization and
// uses the gin.Context argument to send the SVG data back to the user's browser.
func scoreFunction(c *gin.Context) {
	var err error
	cacheKey := "score"
	size := 0
	sizeStr := c.Query("size")
	slog.Info("scoreFunction", "size", sizeStr)
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			slog.Warn("size URL arg -> int", "err", err)
			size = 0
		}
		if size < 0 {
			slog.Warn("size URL arg too low", "size", size, "err", err)
			size = 0
		}
		if size > 0 {
			high := len(sizes) - 1
			if size > high {
				slog.Warn("size URL arg too high", "size", size, "high", high, "err", err)
				size = high
			}
		}
	}
	cacheKey = cacheKey + ":" + strconv.Itoa(size)
	chartCacheMutex.Lock()
	ch, found := chartCache[cacheKey]
	chartCacheMutex.Unlock()
	if !found {
		graph := scoreChart(sizes[size])
		var buf bytes.Buffer
		if err := graph.Render(chart.SVG, &buf); err != nil {
			log.Println(err.Error())
		} else {
			ch = buf.Bytes()
			chartCacheMutex.Lock()
			chartCache[cacheKey] = ch
			chartCacheMutex.Unlock()
		}
	}
	c.Data(http.StatusOK, chart.ContentTypeSVG, ch)
}

// -----------------------------------------------------------------------------

const (
	// height of chart in pixels
	height = 600

	// width of chart in pixes
	width = 750
)

var (
	insertColor = drawing.Color{
		R: 0xff,
		G: 0x00,
		B: 0x7f,
		A: 0x7f,
	}
	annotationColor = drawing.Color{
		R: 0x3f,
		G: 0x7f,
		B: 0xff,
		A: 0xff,
	}
)

// scoreChart generates a chart.Chart object which is a scatter plot of
// handler benchmark vs. warning scores.
func scoreChart(size *sizeData) chart.Chart {
	slog.Info("scoreChart", "size", size)
	handlers := make(map[data.HandlerTag]*handlerCoords)
	for _, hdlr := range bench.HandlerTags() {
		// Only make handler record if y value is within bounds (above size.low.y).
		if score.HandlerBenchScores(hdlr).Overall >= size.low.y {
			handlers[hdlr] = &handlerCoords{y: score.HandlerBenchScores(hdlr).Overall}
		}
	}
	for _, hdlr := range warns.HandlerTags() {
		// Only add value if there is already a benchmark score.
		if coords, found := handlers[hdlr]; found {
			// Only add x-value if it is within bounds (above size.low.x).
			if score.HandlerWarningScore(hdlr) >= size.low.x {
				coords.x = score.HandlerWarningScore(hdlr)
			} else {
				// The x-value is out of bounds but y-value was in bounds,
				// remove handler record previously added.
				delete(handlers, hdlr)
			}
		}
	}
	aValues := make([]chart.Value2, 0, len(handlers)+1)
	xValues := make([]float64, 0, len(handlers)+1)
	yValues := make([]float64, 0, len(handlers)+1)
	for hdlr, coords := range handlers {
		if coords.y > 0.00001 {
			aValues = append(aValues, chart.Value2{
				Label:  warns.HandlerName(hdlr),
				XValue: coords.x, //  + 1,
				YValue: coords.y,
			})
			xValues = append(xValues, coords.x)
			yValues = append(yValues, coords.y)
		}
	}
	ticks := func(low float64) []chart.Tick {
		result := make([]chart.Tick, 0, 11)
		for t := low; t <= 100.0; {
			result = append(result, chart.Tick{
				Value: t,
				Label: strconv.FormatFloat(t, 'f', 1, 64),
			})
			if math.Remainder(t, 10.0) == 0 {
				t += 10.0
			} else {
				t = math.Ceil(t/10.0) * 10.0
			}
		}
		return result
	}
	series := make([]chart.Series, 0, 5)
	for _, s := range sizes {
		if s.low.x > size.low.x && s.low.y > size.low.y {
			series = append(series, chart.ContinuousSeries{
				Style: chart.Style{
					DotWidth:    chart.Disabled,
					StrokeColor: insertColor,
					StrokeWidth: 1,
				},
				XValues: []float64{s.low.x, s.low.x, 100.0},
				YValues: []float64{100.0, s.low.y, s.low.y},
			})
		}
	}
	return chart.Chart{
		Height: height,
		Width:  width,
		XAxis: chart.XAxis{
			Name:  "Warning Score",
			Range: &chart.ContinuousRange{Min: size.low.x, Max: 100.0, Domain: 100.0},
			Ticks: ticks(size.low.x),
		},
		YAxis: chart.YAxis{
			Name: "Benchmark Score",
			//AxisType: chart.YAxisSecondary, // cuts off axis labels on left
			Range: &chart.ContinuousRange{Min: 0, Max: 100.0, Domain: 100.0},
			Ticks: ticks(size.low.y),
		},
		Series: append(series,
			chart.ContinuousSeries{
				Style: chart.Style{
					DotWidth:         5,
					DotColorProvider: scoreChartColorFunction,
					StrokeWidth:      chart.Disabled,
				},
				XValues: xValues,
				YValues: yValues,
			},
			chart.AnnotationSeries{
				Style: chart.Style{
					DotWidth:    chart.Disabled,
					StrokeColor: annotationColor,
					StrokeWidth: 1,
				},
				Annotations: scoreChartAdjustLabels(aValues, size.adjust),
			}),
	}
}

// -----------------------------------------------------------------------------

var defaultLabelSize = &handlerCoords{
	// Approximate width of a label along the x-axis in percentage coordinates (not pixels).
	// This is a trial-and-error value :-(.
	x: 25.0,

	// Approximate height of a label along the y-axis in percentage coordinates (not pixels).
	// This is a trial-and-error value :-(.
	y: 4.5,
}

// scoreChartAdjustLabels
func scoreChartAdjustLabels(locations []chart.Value2, labelSize *handlerCoords) []chart.Value2 {
	groups := make([]*labelGroup, 0, 10)
	singles := make([]int, 0, len(locations))

	// Pull out groups and collect remaining singles.
location: // OMG I can't believe he's using that named loop thingy!!!
	for i := 0; i < len(locations); i++ {
		loc := locations[i]
		for _, group := range groups {
			if group.contains(i) {
				continue location
			}
		}
		for _, group := range groups {
			if group.overlaps(loc, labelSize) {
				group.add(i, loc)
				continue location
			}
		}
		// Here is where we go O(n-squared).
		// Thankfully our array size will likely never be that big.
		for j := i + 1; j < len(locations); j++ {
			other := locations[j]
			if math.Abs(other.XValue-loc.XValue) < labelSize.x && math.Abs(other.YValue-loc.YValue) < labelSize.y {
				group := newLabelGroup()
				group.add(i, loc)
				group.add(j, other)
				groups = append(groups, group)
				continue location
			}
		}
		singles = append(singles, i)
	}

	// Check the remaining singles against the existing groups.
	for _, index := range singles {
		for _, group := range groups {
			if group.overlaps(locations[index], labelSize) {
				group.add(index, locations[index])
			}
		}
	}
	// Adjust label locations in each group.
	// Singles don't need any adjustment.
	for _, group := range groups {
		group.adjust(locations, labelSize)
	}
	// Return for convenience, likely already edited in place.
	return locations
}

const bigByte = float64(0xFF)

// scoreChartColorFunction returns a color dependent on the chart location.
// The algorithm is intended to make it red in the lower left and green in the upper right
// with gradual color changes in between the two corners
// calculated by measuring the distance from each of the two corners.
// The color pattern should be quarter-circular around each corner and
// muddy brownish along a line from the upper left to the lower right.
func scoreChartColorFunction(xr, yr chart.Range, _ int, x, y float64) drawing.Color {
	ratio := scoreChartRatio(xr.GetMin(), yr.GetMin(), xr.GetMax(), yr.GetMax())
	diagonal := scoreChartDistance(xr.GetMin(), yr.GetMin(), xr.GetMax(), yr.GetMax(), ratio)
	distLow := scoreChartDistance(x, y, xr.GetMin(), yr.GetMin(), ratio)
	distHigh := scoreChartDistance(xr.GetMax(), yr.GetMax(), x, y, ratio)
	return drawing.Color{
		R: scoreChartColor(distLow, diagonal),
		G: scoreChartColor(distHigh, diagonal),
		B: 0x00,
		A: 0xff,
	}
}

// scoreChartRatio is a convenience function for use in scoreChartColorFunction.
func scoreChartRatio(xMin, yMin, xMax, yMax float64) float64 {
	return (xMax - xMin) / (yMax - yMin)
}

// scoreChartColor is a convenience function for use in scoreChartColorFunction.
func scoreChartColor(distance, diagonal float64) uint8 {
	if distance > diagonal {
		return 0
	}
	return byte(bigByte - bigByte*distance/diagonal)
}

// scoreChartDistance is a convenience function for use in scoreChartColorFunction.
func scoreChartDistance(xMin, yMin, xMax, yMax, ratio float64) float64 {
	return math.Sqrt(math.Pow((xMax-xMin)/ratio, 2) + math.Pow(yMax-yMin, 2))
}

// -----------------------------------------------------------------------------

// labelGroup collects labels that are close together into a group and
// implements functionality for adjusting labels along (only) the y-axis.
type labelGroup struct {
	// "Set" of index integers in the group for quick lookup.
	indexMap map[int]bool

	// "List" of index integers in the group for array access.
	indices []int

	// High and low x-coordinate values for all labels in the group.
	xHigh, xLow float64

	// High and low y-coordinate values for all labels in the group.
	yHigh, yLow float64
}

// newLabelGroup returns a new labelGroup.
func newLabelGroup() *labelGroup {
	return &labelGroup{
		indexMap: make(map[int]bool),
		indices:  make([]int, 0, 7),
		// Range is 0..100.0 so these should be just as good as math.MaxFloat64
		// (and not look like alien messaging in a debugging window).
		xLow: 1000.0,
		yLow: 1000.0,
	}
}

// add a label index and chart value to the group.
func (lg *labelGroup) add(index int, loc chart.Value2) {
	if _, found := lg.indexMap[index]; !found {
		lg.indexMap[index] = true
		lg.indices = append(lg.indices, index)
		if loc.XValue < lg.xLow {
			lg.xLow = loc.XValue
		} else if loc.XValue > lg.xHigh {
			lg.xHigh = loc.XValue
		}
		if loc.YValue < lg.yLow {
			lg.yLow = loc.YValue
		} else if loc.YValue > lg.yHigh {
			lg.yHigh = loc.YValue
		}
	}
}

// adjust labels in the group by editing their y-coordinate values until they don't overlap.
// Since labels are much longer in the x-axis this seems reasonable.
func (lg *labelGroup) adjust(locations []chart.Value2, labelSize *handlerCoords) {
	labels := make([]string, 0, len(lg.indices))
	for _, index := range lg.indices {
		labels = append(labels, locations[index].Label)
	}
	if len(locations) < 2 {
		return
	}

	// Sort the indices array based on
	// the y-value of the location indexed by the indices array cell.
	sort.Slice(lg.indices, func(i, j int) bool {
		return locations[lg.indices[i]].YValue <= locations[lg.indices[j]].YValue
	})

	// Initialization of various state variables.
	var upDebt, downDebt float64
	var upIndex, downIndex int
	var center int
	if len(lg.indices)%2 == 0 {
		// Even number of indices.
		center = len(lg.indices) / 2
		downIndex = center - 1
		upIndex = center
		overlap := labelSize.y - (locations[lg.indices[upIndex]].YValue - locations[lg.indices[downIndex]].YValue)
		if overlap > 0 {
			downDebt = overlap / 2.0
			locations[lg.indices[downIndex]].YValue -= downDebt
			upDebt += overlap / 2.0
			locations[lg.indices[upIndex]].YValue += upDebt
		}
	} else {
		// Odd number of indices.
		center = (len(lg.indices) - 1) / 2
		downIndex = center
		upIndex = center
	}

	// Adjust labels in the negative y direction.
	for downIndex > 0 {
		overlap := labelSize.y - (locations[lg.indices[downIndex]].YValue - locations[lg.indices[downIndex-1]].YValue)
		if overlap > 0 {
			downDebt += overlap
			locations[lg.indices[downIndex-1]].YValue -= downDebt
		}
		downIndex--
	}

	// Adjust labels in the positive y direction.
	for upIndex < len(lg.indices)-1 {
		overlap := labelSize.y - (locations[lg.indices[upIndex+1]].YValue - locations[lg.indices[upIndex]].YValue)
		if overlap > 0 {
			upDebt += overlap
			locations[lg.indices[upIndex+1]].YValue += upDebt
		}
		upIndex++
	}
}

// contains returns true if the specified index is part of the group.
func (lg *labelGroup) contains(index int) bool {
	_, found := lg.indexMap[index]
	return found
}

// overlaps returns true if the specified chart value overlaps the group.
func (lg *labelGroup) overlaps(loc chart.Value2, labelSize *handlerCoords) bool {
	return loc.XValue >= lg.xLow-labelSize.x && loc.XValue <= lg.xHigh+labelSize.x &&
		loc.YValue >= lg.yLow-labelSize.y && loc.YValue <= lg.xLow+labelSize.y
}
