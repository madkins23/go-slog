package chart

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
	"github.com/madkins23/go-slog/internal/scoring/score"
)

// -----------------------------------------------------------------------------

const (
	// height of chart in pixels
	height = 600

	// width of chart in pixes
	width = 800
)

var (
	annotationColor = drawing.Color{
		R: 0x3f,
		G: 0x7f,
		B: 0xff,
		A: 0xff,
	}
	insertColor = drawing.Color{
		R: 0xff,
		G: 0x00,
		B: 0x7f,
		A: 0x3f,
	}
	lineColor = drawing.Color{
		R: 0x3f,
		G: 0x3f,
		B: 0x3f,
		A: 0x3f,
	}
)

// -----------------------------------------------------------------------------

type handlerCoords struct {
	x, y score.Value
}

func (hc *handlerCoords) adjust(by score.Value) *handlerCoords {
	return &handlerCoords{
		x: hc.x * by,
		y: hc.y * by,
	}
}

var defaultLabelSize = &handlerCoords{
	// Approximate width of a label along the x-axis in percentage coordinates (not pixels).
	// This is a trial-and-error value :-(.
	x: 20.0,

	// Approximate height of a label along the y-axis in percentage coordinates (not pixels).
	// This is a trial-and-error value :-(.
	y: 4.5,
}

type sizeData struct {
	name   string
	low    handlerCoords
	adjust *handlerCoords
	labelX float64
	sizeX  float64
}

var sizes = []*sizeData{
	{
		name:   "Full Size",
		adjust: defaultLabelSize,
		labelX: 110,
	},
	{
		name: "3/4",
		low: handlerCoords{
			x: 25,
			y: 25,
		},
		adjust: defaultLabelSize.adjust(0.7),
		labelX: 107.5,
		sizeX:  85.0,
	},
	{
		name: "Half",
		low: handlerCoords{
			x: 50,
			y: 50,
		},
		adjust: defaultLabelSize.adjust(0.5),
		labelX: 105,
		sizeX:  90.0,
	},
	{
		name: "Quarter",
		low: handlerCoords{
			x: 75,
			y: 75,
		},
		adjust: defaultLabelSize.adjust(0.4),
		labelX: 102.5,
		sizeX:  93.0,
	},
}

// Score generates an SVG chart for the score visualization and
// uses the gin.Context argument to send the SVG data back to the user's browser.
func Score(c *gin.Context, warns *data.Warnings) {
	size := scoreChartSize(c)
	cacheKey := "score:" + c.Param("keeper")
	cacheKey = cacheKey + ":" + strconv.Itoa(size)
	CacheMutex.Lock()
	ch, found := Cache[cacheKey]
	CacheMutex.Unlock()
	keeper := score.GetKeeper(score.KeeperTag(c.Param("keeper")))
	chartData := newChartData(keeper, warns, sizes[size])
	if !found {
		graph := chartData.generate()
		var buf bytes.Buffer
		if err := graph.Render(chart.SVG, &buf); err != nil {
			log.Println(err.Error())
		} else {
			ch = buf.Bytes()
			CacheMutex.Lock()
			Cache[cacheKey] = ch
			CacheMutex.Unlock()
		}
	}
	c.Data(http.StatusOK, chart.ContentTypeSVG, ch)
}

// scoreChartSize determines the chart size from the gin.Context object.
func scoreChartSize(c *gin.Context) int {
	size := 0
	sizeStr := c.Param("size")
	if sizeStr != "" {
		var err error
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
	return size
}

// -----------------------------------------------------------------------------

type scoreChartDataItem struct {
	label  string
	x, y   float64
	labelY float64
}

type scoreChartData struct {
	handlers map[data.HandlerTag]*handlerCoords
	keeper   *score.Keeper
	items    []scoreChartDataItem
	series   []chart.Series
	size     *sizeData
	warns    *data.Warnings
}

func newChartData(keeper *score.Keeper, warns *data.Warnings, size *sizeData) *scoreChartData {
	cd := &scoreChartData{
		handlers: make(map[data.HandlerTag]*handlerCoords),
		keeper:   keeper,
		series:   make([]chart.Series, 0, 6),
		size:     size,
		warns:    warns,
	}
	// Initializes a map from data.HandlerTag to *handlerCoords.
	// Handler records are only used if the handler scores are within the size of the chart.
	for _, hdlr := range cd.keeper.HandlerTags() {
		// Only make handler record if y value is within bounds (above size.low.y).
		if cd.keeper.Y().ScoreFor(hdlr) >= cd.size.low.y {
			cd.handlers[hdlr] = &handlerCoords{y: cd.keeper.Y().ScoreFor(hdlr)}
		}
	}
	for _, hdlr := range cd.warns.HandlerTags() {
		// Only add value if there is already a benchmark score.
		if coords, found := cd.handlers[hdlr]; found {
			// Only add x-value if it is within bounds (above size.low.x).
			if cd.keeper.X().ScoreFor(hdlr) >= cd.size.low.x {
				coords.x = cd.keeper.X().ScoreFor(hdlr)
			} else {
				// The x-value is out of bounds but y-value was in bounds,
				// remove handler record previously added.
				delete(cd.handlers, hdlr)
			}
		}
	}
	cd.items = make([]scoreChartDataItem, 0, len(cd.handlers)+1)
	return cd
}

func SmallestChartSize(keeper *score.Keeper) uint8 {
	smallest := 0
	for sz := 1; sz < len(sizes); sz++ {
		for _, hdlr := range keeper.HandlerTags() {
			if keeper.Y().ScoreFor(hdlr) < sizes[sz].low.y ||
				keeper.X().ScoreFor(hdlr) < sizes[sz].low.x {
				return uint8(smallest)
			}
		}
		smallest = sz
	}
	return uint8(smallest)
}

// generate a chart.Chart object which is
// a scatter plot of handler benchmark vs. warning scores.
func (cd *scoreChartData) generate() chart.Chart {
	cd.dataMarkers()
	cd.sizeMarkers()
	return chart.Chart{
		Height: height,
		Width:  width,
		XAxis: chart.XAxis{
			Name:  cd.keeper.X().Name(),
			Range: &chart.ContinuousRange{Min: float64(cd.size.low.x), Max: 100.0, Domain: 100.0},
			Ticks: scoreChartTicks(float64(cd.size.low.x)),
		},
		YAxis: chart.YAxis{
			Name: cd.keeper.Y().Name(),
			//AxisType: chart.YAxisSecondary, // cuts off axis labels on left
			Range: &chart.ContinuousRange{Min: 0, Max: 100.0, Domain: 100.0},
			Ticks: scoreChartTicks(float64(cd.size.low.y)),
		},
		Series: append(cd.series,
			chart.ContinuousSeries{
				Style: chart.Style{
					DotWidth:         5,
					DotColorProvider: scoreChartColorFunction,
					StrokeWidth:      chart.Disabled,
				},
				XValues: cd.xValues(),
				YValues: cd.yValues(),
			},
			chart.AnnotationSeries{
				Style: chart.Style{
					DotWidth:    chart.Disabled,
					StrokeColor: annotationColor,
					StrokeWidth: 1,
				},
				Annotations: cd.annotations(),
			}),
	}
}

// annotations returns the label items for the data items.
func (cd *scoreChartData) annotations() []chart.Value2 {
	result := make([]chart.Value2, len(cd.items))
	for i, item := range cd.items {
		result[i] = chart.Value2{
			Label:  item.label,
			XValue: cd.size.labelX,
			YValue: item.labelY,
		}
	}
	return result
}

// dataMarkers generates chart points and labels to represent score values.
func (cd *scoreChartData) dataMarkers() {
	for hdlr, coords := range cd.handlers {
		if coords.y > 0.00001 {
			x := float64(coords.x)
			y := float64(coords.y)
			cd.items = append(cd.items, scoreChartDataItem{
				label: cd.warns.HandlerName(hdlr),
				x:     x,
				y:     y,
			})
		}
	}
	necessary := float64(len(cd.handlers))*(float64(cd.size.adjust.y)+1) + 1
	// Get the vertical range of the data items.
	bottom := 100.0
	top := 0.0
	for _, item := range cd.items {
		if item.y > top {
			top = item.y
		}
		if item.y < bottom {
			bottom = item.y
		}
	}
	top = (top+bottom)/2 + necessary/2
	if top > 100 {
		top = 100
	}
	// Sort data items by decreasing y-value.
	sort.Slice(cd.items, func(i, j int) bool {
		return cd.items[i].y >= cd.items[j].y
	})
	for i, item := range cd.items {
		cd.series = append(cd.series, chart.ContinuousSeries{
			Style: chart.Style{
				DotWidth:    chart.Disabled,
				StrokeColor: lineColor,
				StrokeWidth: 1,
			},
			XValues: []float64{item.x, cd.size.labelX},
			YValues: []float64{item.y, top},
		})
		cd.items[i].labelY = top
		top -= float64(cd.size.adjust.y) + 1
	}
}

// sizeMarkers generates chart lines and labels to represent chart size options.
func (cd *scoreChartData) sizeMarkers() {
	var labelX float64
	labels := make([]chart.Value2, 0, 3)
	for _, s := range sizes {
		if s.low.x > cd.size.low.x && s.low.y > cd.size.low.y {
			if labelX == 0.0 {
				labelX = s.sizeX
			}
			lowX := float64(s.low.x)
			lowY := float64(s.low.y)
			cd.series = append(cd.series, chart.ContinuousSeries{
				Style: chart.Style{
					DotWidth:    chart.Disabled,
					StrokeColor: insertColor,
					StrokeWidth: 1,
				},
				XValues: []float64{lowX, lowX, 100.0},
				YValues: []float64{100.0, lowY, lowY},
			})
			labels = append(labels, chart.Value2{
				Label:  s.name,
				XValue: labelX,
				YValue: lowY,
			})
		}
	}
	if len(labels) > 0 {
		cd.series = append(cd.series, chart.AnnotationSeries{
			Style: chart.Style{
				DotWidth:    chart.Disabled,
				StrokeColor: insertColor,
				StrokeWidth: 1,
			},
			Annotations: labels,
		})
	}
}

func (cd *scoreChartData) xValues() []float64 {
	result := make([]float64, len(cd.items))
	for i, item := range cd.items {
		result[i] = item.x
	}
	return result
}

func (cd *scoreChartData) yValues() []float64 {
	result := make([]float64, len(cd.items))
	for i, item := range cd.items {
		result[i] = item.y
	}
	return result
}

// -----------------------------------------------------------------------------

// scoreChartColor is a convenience function for use in scoreChartColorFunction.
func scoreChartColor(distance, diagonal float64) uint8 {
	if distance > diagonal {
		return 0
	}
	return byte(math.MaxUint8 - math.MaxUint8*distance/diagonal)
}

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

// scoreChartDistance is a convenience function for use in scoreChartColorFunction.
func scoreChartDistance(xMin, yMin, xMax, yMax, ratio float64) float64 {
	return math.Sqrt(math.Pow((xMax-xMin)/ratio, 2) + math.Pow(yMax-yMin, 2))
}

// scoreChartRatio is a convenience function for use in scoreChartColorFunction.
func scoreChartRatio(xMin, yMin, xMax, yMax float64) float64 {
	return (xMax - xMin) / (yMax - yMin)
}

// scoreChartTicks returns a list of chart.Tick objects for an axis.
func scoreChartTicks(low float64) []chart.Tick {
	tickDistance := 10.0
	if 100.0-low <= 10.0 {
		tickDistance = 2.0
	} else if 100.0-low <= 25.0 {
		tickDistance = 2.0
	} else if 100.0-low <= 50.0 {
		tickDistance = 5.0
	}
	result := make([]chart.Tick, 0, 11)
	for t := low; t <= 100.0; {
		result = append(result, chart.Tick{
			Value: t,
			Label: strconv.FormatFloat(t, 'f', 0, 64),
		})
		if math.Remainder(t, tickDistance) == 0 {
			t += tickDistance
		} else {
			t = math.Ceil(t/tickDistance) * tickDistance
		}
	}
	return result
}
