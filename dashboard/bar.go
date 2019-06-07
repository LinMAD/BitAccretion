package dashboard

import (
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/barchart"
)

// BarWidgetHandler for dashboard
type BarWidgetHandler struct {
	name                  string
	isOkRequestsToCollect bool
	barChart              *barchart.BarChart
	// barChartMap mapping relation between graph and bar chart, order is sensitive
	barChartMap map[string]*model.SystemMetric
}

// HandleNotifyEvent update bar chat data
func (bw *BarWidgetHandler) HandleNotifyEvent(e event.UpdateEvent) error {
	var max int

	vertices := e.MonitoringGraph.GetAllVertices()
	l := len(vertices)

	if bw.isOkRequestsToCollect {
		max = bw.getMaxRequestValue(false, &e.MonitoringGraph)
	} else {
		max = bw.getMaxRequestValue(true, &e.MonitoringGraph)
	}

	// Set new data from metrics in same order as before
	for i := 0; i < l; i++ {
		if bw.isOkRequestsToCollect {
			bw.barChartMap[vertices[i].Name] = &vertices[i].Metric
		} else {
			bw.barChartMap[vertices[i].Name] = &vertices[i].Metric
		}
	}

	i := 0
	bar := make([]int, l)
	for _, m := range bw.barChartMap {
		if bw.isOkRequestsToCollect {
			bar[i] = int(m.RequestCount)
		} else {
			bar[i] = int(m.ErrorCount)
		}
		i++
	}

	return bw.barChart.Values(bar, max)
}

// GetName of widget handler
func (bw *BarWidgetHandler) GetName() string {
	return bw.name
}

// getMaxRequestValue max value from vertices
func (bw *BarWidgetHandler) getMaxRequestValue(isErrorsReqs bool, g *model.Graph) int {
	max := 1
	allVertices := g.GetAllVertices()

	for i := 0; i < len(allVertices); i++ {
		if isErrorsReqs && max < int(allVertices[i].Metric.ErrorCount) {
			max = int(allVertices[i].Metric.ErrorCount)
			continue
		}

		if max < int(allVertices[i].Metric.RequestCount) {
			max = int(allVertices[i].Metric.RequestCount)
		}
	}

	return max
}

// NewBarWidget creates and returns prepared widget
func NewBarWidget(name string, barColor cell.Color, isOkReqs bool, nodes []*model.Node) (*BarWidgetHandler, error) {
	sysCount := len(nodes)
	sysNames := make([]string, sysCount)
	sysBarColors := make([]cell.Color, sysCount)
	sysBarValColors := make([]cell.Color, sysCount)
	sysBarValMap := make(map[string]*model.SystemMetric)

	for i := 0; i < sysCount; i++ {
		sysBarColors[i] = barColor
		sysBarValColors[i] = cell.ColorWhite
		sysBarValMap[nodes[i].Name] = nil
	}

	i := 0
	for n := range sysBarValMap {
		sysNames[i] = n
		i++
	}

	sysBar, sysBarErr := barchart.New(
		barchart.BarColors(sysBarColors),
		barchart.ValueColors(sysBarValColors),
		barchart.ShowValues(),
		barchart.Labels(sysNames),
	)
	if sysBarErr != nil {
		return nil, sysBarErr
	}

	widget := &BarWidgetHandler{
		name:                  name,
		barChart:              sysBar,
		isOkRequestsToCollect: isOkReqs,
		barChartMap:           sysBarValMap,
	}

	return widget, nil
}
