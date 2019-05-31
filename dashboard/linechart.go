package dashboard

import (
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/linechart"
)

// SparkLineWidgetHandler for dashboard
type SparkLineWidgetHandler struct {
	name     string
	lc *linechart.LineChart
}

// HandleNotifyEvent update spark line chat data
func (s SparkLineWidgetHandler) HandleNotifyEvent(e event.UpdateEvent) {
	vertices := e.MonitoringGraph.GetAllVertices()
	l := len(vertices)

	sData := make([]float64, l)
	eData := make([]float64, l)
	for i := 0; i < l; i++ {
		sData[i] = float64(vertices[i].Metric.RequestCount)
		eData[i] = float64(vertices[i].Metric.ErrorCount)
	}

	lcErr := s.lc.Series(
		"ok",
		sData,
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorGreen)),
	)
	if lcErr != nil {
		panic(lcErr)
	}
	if err := s.lc.Series("bad", eData, linechart.SeriesCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
}

// GetName of widget handler
func (s SparkLineWidgetHandler) GetName() string {
	return s.name
}

// NewSparkLineChart creates and returns prepared widget
func NewSparkLineChart(color cell.Color, name string, nodes []model.Node) (*SparkLineWidgetHandler, error) {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorWhite)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorWhite)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorWhite)),
	)
	if err != nil {
		return nil, err
	}

	widget := &SparkLineWidgetHandler{
		name:     name,
		lc: lc,
	}

	return widget, nil
}
