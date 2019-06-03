package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/LinMAD/BitAccretion/dashboard"
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/logger"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/LinMAD/BitAccretion/stub"

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

	// TODO Get data from provider
	graph := stub.GetStubGraph()

	// TODO Use same context to cancel all widgets subscribers
	ctx, cancel := context.WithCancel(context.Background())
	monitoringDashboard, err := dashboard.NewMonitoringDashboard("BitAccretion", t, graph)
	if err != nil {
		panic(err)
	}

	// TODO used ctx to cancel all widgets
	monitoringObserver := event.NewDashboardObserver(monitoringDashboard.EventLogger)
	monitoringObserver.RegisterNewSubscriber(monitoringDashboard)

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	// TODO Replace that
	go playMonitoring(monitoringObserver, monitoringDashboard.EventLogger, 1*time.Second)

	runErr := termdash.Run(
		ctx,
		t,
		monitoringDashboard.TerminalContainer,
		termdash.KeyboardSubscriber(quitter),
		termdash.RedrawInterval(1*time.Second),
	)
	if runErr != nil {
		panic(err)
	}
}

// TODO Remove stub and make with channels
func playMonitoring(monitoringObserver event.IObserver, log logger.ILogger, delay time.Duration) []model.Node {
	nodes := stub.GetStubNodes()

	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Normal("Received new monitoring graph update")

			graph := model.NewGraph()
			for i := 0; i < len(nodes); i++ {
				nodes[i].Metric.RequestCount = rand.Intn(15000)
				nodes[i].Metric.ErrorCount = rand.Intn(15000)

				graph.AddVertex(model.VertexName(nodes[i].Name), nodes[i])
			}

			monitoringObserver.NotifySubscribers(event.UpdateEvent{
				MonitoringGraph: *graph,
			})
		}
	}
}
