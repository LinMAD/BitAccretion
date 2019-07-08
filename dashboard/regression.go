package dashboard

import (
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/gauge"
)

// GaugeRegressHandler for dashboard
type GaugeRegressHandler struct {
	name   string
	gauge  *gauge.Gauge
	config *model.Config
	// currentRatio showing regression of errors
	currentRatio int
	prevErrCount int
}

// HandleNotifyEvent update of gauge regression process
func (g *GaugeRegressHandler) HandleNotifyEvent(e event.UpdateEvent) error {
	if err := g.gauge.Percent(g.currentRatio); err != nil {
		panic(err)
	}

	var newErrCount int
	vertices := e.MonitoringGraph.GetAllVertices()
	l := len(vertices)

	for i := 0; i < l; i++ {
		newErrCount += int(vertices[i].Metric.ErrorCount)
	}

	// TODO Think about longer compactions maybe between hours or weeks
	if newErrCount == 0 {
		g.currentRatio = 0
		g.prevErrCount = 0
	} else {
		g.currentRatio = +(newErrCount - g.prevErrCount)
		if g.currentRatio > 100 {
			g.currentRatio = 100
		} else if g.currentRatio < 0 {
			g.currentRatio = 0
		}
	}

	g.prevErrCount = newErrCount

	return nil
}

// GetName of widget handler
func (g *GaugeRegressHandler) GetName() string {
	return g.name
}

// NewRegressionWidget create and return prepared widget, shows % of ok/err requests as regression
func NewRegressionWidget(name string, c *model.Config) (*GaugeRegressHandler, error) {
	g, err := gauge.New(
		gauge.Height(2),
		gauge.Color(cell.ColorRed),
		gauge.Border(linestyle.Light, cell.FgColor(cell.ColorYellow)),
		gauge.BorderTitle(name),
	)
	if err != nil {
		return nil, err
	}

	widget := &GaugeRegressHandler{
		name:   name,
		gauge:  g,
		config: c,
	}

	return widget, nil
}
