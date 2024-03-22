package main

import (
	"bytes"
	"log"
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

// scoreFunction generates an SVG chart for the score visualization and
// uses the gin.Context argument to send the SVG data back to the user's browser.
func scoreFunction(c *gin.Context) {
	cacheKey := "score"
	chartCacheMutex.Lock()
	ch, found := chartCache[cacheKey]
	chartCacheMutex.Unlock()
	if !found {
		graph := scoreChart()
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

// scoreChart generates a chart.Chart object which is a scatter plot of
// handler benchmark vs. warning scores.
func scoreChart() chart.Chart {
	type handlerCoords struct{ x, y float64 }
	handlers := make(map[data.HandlerTag]*handlerCoords)
	for _, hdlr := range bench.HandlerTags() {
		handlers[hdlr] = &handlerCoords{y: score.HandlerBenchScores(hdlr).Overall}
	}
	for _, hdlr := range warns.HandlerTags() {
		if coords, found := handlers[hdlr]; found {
			// Only add value if there is already a benchmark score.
			coords.x = score.HandlerWarningScore(hdlr)
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
	ticks := make([]chart.Tick, 11)
	for i := 0; i < 11; i++ {
		ticks[i].Value = float64(i * 10)
		ticks[i].Label = strconv.FormatFloat(ticks[i].Value, 'f', 1, 64)
	}
	return chart.Chart{
		Height: height,
		Width:  width,
		XAxis: chart.XAxis{
			Name:  "Warning Score",
			Range: &chart.ContinuousRange{Min: 0, Max: 100.0, Domain: 100.0},
			Ticks: ticks,
		},
		YAxis: chart.YAxis{
			Name: "Benchmark Score",
			//AxisType: chart.YAxisSecondary, // cuts off axis labels on left
			Range: &chart.ContinuousRange{Min: 0, Max: 100.0, Domain: 100.0},
			Ticks: ticks,
		},
		Series: []chart.Series{
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
					StrokeWidth: chart.Disabled,
					DotWidth:    chart.Disabled,
				},
				Annotations: scoreChartAdjustLabels(aValues),
			},
		},
	}
}

// -----------------------------------------------------------------------------

// scoreChartAdjustLabels
func scoreChartAdjustLabels(locations []chart.Value2) []chart.Value2 {
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
			if group.overlaps(loc) {
				group.add(i, loc)
				continue location
			}
		}
		for j := i + 1; j < len(locations); j++ {
			other := locations[j]
			if math.Abs(other.XValue-loc.XValue) < labelWidth && math.Abs(other.YValue-loc.YValue) < labelHeight {
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
			if group.overlaps(locations[index]) {
				group.add(index, locations[index])
			}
		}
	}
	// Adjust label locations in each group.
	// Singles don't need any adjustment.
	for _, group := range groups {
		group.adjust(locations)
	}
	// Return for convenience, likely already edited in place.
	return locations
}

// -----------------------------------------------------------------------------

const bigByte = float64(0xFF)

// scoreChartColorFunction returns a color dependent on the chart location.
// The algorithm is intended to make it red in the lower left and green in the upper right
// with gradual color changes in between.
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

// labelGroup collects labels that are close together into a group.
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

const (
	// Approximate width of a label along the x-axis in percentage coordinates (not pixels).
	// This is a trial-and-error value :-(.
	labelWidth = 25.0

	// Approximate height of a label along the y-axis in percentage coordinates (not pixels).
	// This is a trial-and-error value :-(.
	labelHeight = 4.5
)

// adjust labels in the group by editing their y-coordinate values until they don't overlap.
func (lg *labelGroup) adjust(locations []chart.Value2) {
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
		overlap := labelHeight - (locations[lg.indices[upIndex]].YValue - locations[lg.indices[downIndex]].YValue)
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
		overlap := labelHeight - (locations[lg.indices[downIndex]].YValue - locations[lg.indices[downIndex-1]].YValue)
		if overlap > 0 {
			downDebt += overlap
			locations[lg.indices[downIndex-1]].YValue -= downDebt
		}
		downIndex--
	}

	// Adjust labels in the positive y direction.
	for upIndex < len(lg.indices)-1 {
		overlap := labelHeight - (locations[lg.indices[upIndex+1]].YValue - locations[lg.indices[upIndex]].YValue)
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
func (lg *labelGroup) overlaps(loc chart.Value2) bool {
	return loc.XValue >= lg.xLow-labelWidth && loc.XValue <= lg.xHigh+labelWidth &&
		loc.YValue >= lg.yLow-labelHeight && loc.YValue <= lg.xLow+labelHeight
}
