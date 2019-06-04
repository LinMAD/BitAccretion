package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"plugin"
	"time"

	"github.com/LinMAD/BitAccretion/dashboard"
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/logger"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/LinMAD/BitAccretion/provider"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

var (
	providerImpl provider.IProvider
	configPath   string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic("Could not retrieve working directory, error: " + err.Error())
	}
	configPath = wd + "/config.json"

	mod, err := plugin.Open(wd + "/provider.so")
	if err != nil {
		panic("Unable to open provider.so plugin, error: " + err.Error())
	}

	// Validate plugin - lookup for exported base function to get implementation
	prc, err := mod.Lookup("NewProvider")
	if err != nil {
		log.Fatalf("Expected to be exported Processor structure in plugin, err: %v", err)
	}

	// Add implemented plugin to kernel
	providerImpl = prc.(func() provider.IProvider)()
}

func main() {
	t, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	// TODO Move that to other place (Plugin warm up)
	/**
	PLUGIN WARM UP
	*/

	// Load plugin settings
	pluginCfgErr := providerImpl.LoadConfig(configPath)
	if pluginCfgErr != nil {
		panic("Provider configuration, error: " + pluginCfgErr.Error())
	}

	// Boot plugin
	pluginBootErr := providerImpl.Boot()
	if pluginBootErr != nil {
		panic("Provider boot, error: " + pluginBootErr.Error())
	}

	// Get first system graph
	var providerGraph model.Graph
	var providerGraphErr error
	providerGraph, providerGraphErr = providerImpl.DispatchMonitoredData()
	if providerGraphErr != nil {
		panic("Provider dispatch monitoring data failed, error: " + providerGraphErr.Error())
	}

	// TODO Move that to other place
	/**
	Dashboard creation
	*/

	// TODO Use same context to cancel all widgets subscribers
	ctx, cancel := context.WithCancel(context.Background())
	monitoringDashboard, err := dashboard.NewMonitoringDashboard("BitAccretion", t, providerGraph)
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

	/**
	Observe data changes of system
	*/
	// TODO Move that to pckg
	go providerMonitoring(monitoringObserver, monitoringDashboard.EventLogger, 1*time.Second)

	/**
	Handle dashboard view
	*/
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

// TODO Remove tmp func
func providerMonitoring(monitoringObserver event.IObserver, log logger.ILogger, delay time.Duration) {
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Normal("Requesting for new monitoring graph update")

			providerGraph, providerGraphErr := providerImpl.DispatchMonitoredData()
			if providerGraphErr != nil {
				log.Error(providerGraphErr.Error())
				return
			}

			for _, n := range providerGraph.GetAllVertices() {
				log.Normal(
					fmt.Sprintf("New metrics of %s OK: %d ERR: %d", n.Name, n.Metric.RequestCount, n.Metric.ErrorCount),
				)
			}

			monitoringObserver.NotifySubscribers(event.UpdateEvent{
				MonitoringGraph: providerGraph,
			})
		}
	}
}
