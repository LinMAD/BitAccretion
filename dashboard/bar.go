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

	data := make([]int, l)
	for i := 0; i < l; i++ {
		if bw.isOkRequestsToCollect {
			data[i] = vertices[i].Metric.RequestCount
		} else {
			data[i] = vertices[i].Metric.ErrorCount
		}
	}

	return bw.barChart.Values(data, max)
}

// GetName of widget handler
func (bw *BarWidgetHandler) GetName() string {
	return bw.name
}

// getMaxRequestValue max value from vertices
func (bw *BarWidgetHandler) getMaxRequestValue(isErrorsReqs bool, g *model.Graph) (max int) {
	allVertices := g.GetAllVertices()

	for i := 0; i < len(allVertices); i++ {
		if isErrorsReqs {
			if max < allVertices[i].Metric.ErrorCount {
				max = allVertices[i].Metric.ErrorCount
			}

			continue
		}

		if max < allVertices[i].Metric.RequestCount {
			max = allVertices[i].Metric.RequestCount
		}
	}

	return
}

// NewBarWidget creates and returns prepared widget
func NewBarWidget(name string, barColor cell.Color, isOkReqs bool, nodes []model.Node) (*BarWidgetHandler, error) {
	sysCount := len(nodes)
	sysNames := make([]string, sysCount)
	sysBarColors := make([]cell.Color, sysCount)
	sysValBarColors := make([]cell.Color, sysCount)

	for i := 0; i < sysCount; i++ {
		sysNames[i] = nodes[i].Name
		sysBarColors[i] = barColor
		sysValBarColors[i] = cell.ColorWhite
	}

	sysBar, sysBarErr := barchart.New(
		barchart.BarColors(sysBarColors),
		barchart.ValueColors(sysValBarColors),
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
	}

	return widget, nil
}
