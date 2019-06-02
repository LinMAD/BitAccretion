package dashboard

import (
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/linechart"
)

// maxPoints in line chart for one line (control visual overflow and data updates)
const maxPoints = 50

// SparkLineWidgetHandler for dashboard
type SparkLineWidgetHandler struct {
	name  string
	lc    *linechart.LineChart
	lines seriesData
}

// seriesData used to draw points in line chart
type seriesData struct {
	okData  []float64
	badData []float64
}

// HandleNotifyEvent update spark line chat data
func (s *SparkLineWidgetHandler) HandleNotifyEvent(e event.UpdateEvent) {
	s.updateLineData(&e.MonitoringGraph)

	okLineErr := s.lc.Series(
		"ok",
		s.lines.okData,
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.SeriesXLabels(map[int]string{0: "Iteration: "}),
	)
	if okLineErr != nil {
		panic(okLineErr) // TODO Handle in grace way, log or ignore
	}

	badLineErr := s.lc.Series("bad", s.lines.badData, linechart.SeriesCellOpts(cell.FgColor(cell.ColorRed)))
	if badLineErr != nil {
		panic(badLineErr) // TODO Handle in grace way, log or ignore
	}
}

// GetName of widget handler
func (s *SparkLineWidgetHandler) GetName() string {
	return s.name
}

// updateLineData in widget with stored history of each prev points
func (s *SparkLineWidgetHandler) updateLineData(g *model.Graph) {
	var okPoints, badPoints []float64
	nodes := g.GetAllVertices()

	if len(s.lines.okData) >= maxPoints {
		okPoints = s.lines.okData[1:maxPoints]
	} else {
		okPoints = s.lines.okData
	}

	if len(s.lines.badData) >= maxPoints {
		badPoints = s.lines.badData[1:maxPoints]
	} else {
		badPoints = s.lines.badData
	}

	var okPoint, badPoint float64
	for i := 0; i < len(nodes); i++ {
		okPoint += float64(nodes[i].Metric.RequestCount)
		badPoint += float64(nodes[i].Metric.ErrorCount)
	}

	s.lines.okData = append(okPoints, okPoint)
	s.lines.badData = append(badPoints, badPoint)
}

// NewLineWidget creates and returns prepared widget
func NewLineWidget(name string) (*SparkLineWidgetHandler, error) {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorWhite)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorWhite)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorWhite)),
	)
	if err != nil {
		return nil, err
	}

	widget := &SparkLineWidgetHandler{
		name: name,
		lc:   lc,
		lines: seriesData{
			okData:  make([]float64, 0),
			badData: make([]float64, 0),
		},
	}

	return widget, nil
}
