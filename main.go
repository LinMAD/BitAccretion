package main

import (
	"context"
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/LinMAD/BitAccretion/stub"
	"math/rand"
	"time"

	"github.com/LinMAD/BitAccretion/dashboard"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

func main() {
	t, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	// TODO Use same context to cancel all widgets subscribers
	ctx, cancel := context.WithCancel(context.Background())

	monitoringDashboard, err := dashboard.NewMonitoringDashboard(t)
	if err != nil {
		panic(err)
	}

	monitoringObserver := event.NewDashboardObserver()
	monitoringObserver.RegisterSubscriber(monitoringDashboard)

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	go playMonitoring(monitoringObserver, 1 * time.Second)

	runErr := termdash.Run(
		ctx,
		t,
		monitoringDashboard.TerminalContainer,
		termdash.KeyboardSubscriber(quitter),
		termdash.RedrawInterval(500*time.Millisecond),
	)
	if runErr != nil {
		panic(err)
	}
}

// TODO Remove stub and make with channels
func playMonitoring(monitoringObserver event.IWidgetObserver, delay time.Duration) []model.Node {
	nodes := stub.GetStubNodes()

	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			graph := model.NewGraph()
			for i := 0; i < len(nodes); i++ {
				nodes[i].Metric.RequestCount = rand.Intn(15000)
				nodes[i].Metric.ErrorCount = rand.Intn(200)

				graph.AddVertex(model.VertexName(nodes[i].Name), nodes[i])
			}

			monitoringObserver.NotifySubscribers(event.UpdateEvent{
				MonitoringGraph: *graph,
			})
		}
	}
}
