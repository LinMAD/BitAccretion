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

	max = bw.getMaxRequestValue(bw.isOkRequestsToCollect, &e.MonitoringGraph)
	vertices := e.MonitoringGraph.GetAllVertices()
	l := len(vertices)

	// Labels must be collected with same order as values
	barLabels := make([]string, l)
	barValues := make([]int, l)
	for i := 0; i < l; i++ {
		vertex := vertices[i]
		barLabels[i] = vertex.Name

		if bw.isOkRequestsToCollect {
			barValues[i] = int(vertex.Metric.RequestCount)
		} else {
			barValues[i] = int(vertex.Metric.ErrorCount)
		}
	}

	// Update bar with new collected data set
	return bw.barChart.Values(barValues, max, barchart.Labels(barLabels))
}

// GetName of widget handler
func (bw *BarWidgetHandler) GetName() string {
	return bw.name
}

// getMaxRequestValue max value from vertices
func (bw *BarWidgetHandler) getMaxRequestValue(isOkReq bool, g *model.Graph) int {
	max := 1
	allVertices := g.GetAllVertices()

	for i := 0; i < len(allVertices); i++ {
		vertex := allVertices[i]

		if isOkReq && max < int(vertex.Metric.RequestCount) {
			max = int(vertex.Metric.RequestCount)
			continue
		}

		if max < int(vertex.Metric.ErrorCount) {
			max = int(vertex.Metric.ErrorCount)
		}
	}

	return max
}

// NewBarWidget creates and returns prepared widget
func NewBarWidget(name string, barColor cell.Color, isOkReqs bool, nodes []*model.Node) (*BarWidgetHandler, error) {
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
