package dashboard

import (
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/barchart"
	"math/rand"
)

// BarchartWidgetHandler for dashboard, shows in bars requests of each system in node
type BarchartWidgetHandler struct {
	name     string
	barChart *barchart.BarChart
}

// HandleNotifyEvent update bar chat data
func (barWidget *BarchartWidgetHandler) HandleNotifyEvent(e event.UpdateEvent) {
	maxReqVal := 0
	allNodes := e.MonitoringGraph.GetAllVertices()
	lenNodes := len(allNodes)

	for i := 0; i < lenNodes; i++ {
		if maxReqVal < allNodes[i].Metric.RequestCount {
			maxReqVal = allNodes[i].Metric.RequestCount
		}
	}

	var barData []int // TODO Try allocate space of slice
	for i := 0; i < lenNodes; i++ {
		barData = append(barData, int(rand.Int31n(int32(maxReqVal))))
	}

	// Add values to bar barChart and put max value of it
	if err := barWidget.barChart.Values(barData, maxReqVal); err != nil {
		panic(err) // TODO Handle in grace way, log or ignore
	}
}

// GetName of widget handler
func (barWidget *BarchartWidgetHandler) GetName() string {
	return barWidget.name
}

// NewBarChart creates and returns prepared dashboard widget
func NewBarChart(name string, nodes []model.Node) (*BarchartWidgetHandler, error) {
	barWidth := 0
	sysCount := len(nodes)
	sysNames := make([]string, sysCount)
	sysBarColors := make([]cell.Color, sysCount)
	sysValBarColors := make([]cell.Color, sysCount)

	// TODO Think is it really need to color on new health, if so must be done in event (Is it possible in API?)
	for i := 0; i < sysCount; i++ {
		sysNames[i] = nodes[i].Name

		switch nodes[i].Health {
		case model.HealthWarning:
			sysBarColors[i] = cell.ColorYellow
			sysValBarColors[i] = cell.ColorWhite
		case model.HealthCritical:
			sysBarColors[i] = cell.ColorRed
			sysValBarColors[i] = cell.ColorWhite
		default:
			sysBarColors[i] = cell.ColorBlue
			sysValBarColors[i] = cell.ColorWhite
		}

		if barWidth < len(nodes[i].Name) {
			barWidth = len(nodes[i].Name)
		}
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

	widget := &BarchartWidgetHandler{
		name:     name,
		barChart: sysBar,
	}

	return widget, nil
}
