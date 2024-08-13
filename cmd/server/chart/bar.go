package chart

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vicanso/go-charts/v2"

	"github.com/madkins23/go-slog/internal/data"
)

// Bar generates an SVG chart for the current object tags.
func Bar(c *gin.Context, bench *data.Benchmarks) {
	itemArg := strings.TrimSuffix(c.Param("item"), ".svg")
	item, err := data.BenchItemsString(itemArg)
	if err != nil {
		slog.Error("Bad URL parameter", "param", itemArg, "err", err)
		// TODO: what to do here?
		return
	}
	tag := c.Param("tag")
	cacheKey := tag + ":" + item.String()
	CacheMutex.Lock()
	ch, found := Cache[cacheKey]
	CacheMutex.Unlock()
	if !found {
		var labels []string
		var values []float64
		if records := bench.HandlerRecordsFor(data.TestTag(tag)); records != nil {
			labels, values = chartTest(bench, records, item)
		} else if records := bench.TestRecordsFor(data.HandlerTag(tag)); records != nil {
			labels, values = chartHandler(bench, records, item)
		} else {
			slog.Error("Neither handler nor benchmark records found", "fn", "chartFunction")
			c.HTML(http.StatusBadRequest, "pageFunction", gin.H{
				"ErrorTitle":   "Template failed execution",
				"ErrorMessage": "No records for " + tag})
			return
		}
		const verticalPadding = 100
		const barMargin = 5
		const barWidth = 15
		painter, err := charts.HorizontalBarRender(
			[][]float64{values},
			charts.SVGTypeOption(),
			charts.TitleTextOptionFunc(item.LongName()),
			charts.YAxisDataOptionFunc(labels),
			charts.WidthOptionFunc(400),
			charts.HeightOptionFunc(verticalPadding+(barMargin+barWidth)*len(values)),
			charts.PaddingOptionFunc(charts.Box{
				Top:    10,
				Right:  20,
				Bottom: 10,
				Left:   10,
			}),
			func(opt *charts.ChartOption) {
				opt.BarWidth = barWidth
			},
		)
		if err != nil {
			panic(err)
		}
		ch, err = painter.Bytes()
		if err != nil {
			panic(err)
		}
		CacheMutex.Lock()
		Cache[cacheKey] = ch
		CacheMutex.Unlock()
	}
	c.Data(http.StatusOK, "image/svg+xml", ch)
}

// charTest returns labels and values for a Test chart.
func chartTest(bench *data.Benchmarks, records data.HandlerRecords, item data.BenchItems) (labels []string, values []float64) {
	labels = make([]string, 0, len(records))
	values = make([]float64, 0, len(records))
	for _, tag := range bench.HandlerTags() {
		if record, found := records[tag]; found {
			labels = append(labels, bench.HandlerName(tag))
			values = append(values, record.ItemValue(item))
		}
	}
	reverse(labels)
	reverse(values)
	return
}

// chartHandler returns labels and values for a Handler chart.
func chartHandler(bench *data.Benchmarks, records data.TestRecords, item data.BenchItems) (labels []string, values []float64) {
	labels = make([]string, 0, len(records))
	values = make([]float64, 0, len(records))
	for _, tag := range bench.TestTags() {
		if record, found := records[tag]; found {
			labels = append(labels, bench.TestName(tag))
			values = append(values, record.ItemValue(item))
		}
	}
	reverse(labels)
	reverse(values)
	return
}

// -----------------------------------------------------------------------------

// reverse an array.
func reverse[T any](array []T) {
	i := 0
	j := len(array) - 1
	for i < j {
		array[i], array[j] = array[j], array[i]
		i++
		j--
	}
}
