package kernel

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/LinMAD/BitAccretion/dashboard"
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/logger"
	"github.com/LinMAD/BitAccretion/provider"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

// Kernel core of whole application it's managing states and data communications
type Kernel struct {
	d *dashboard.MonitoringDashboard
	p provider.IProvider
	o event.IObserver
	l logger.ILogger
}

// NewKernel of monitoring
func NewKernel(dataProvider provider.IProvider) *Kernel {
	k := &Kernel{
		p: dataProvider,
	}

	return k
}

//initProvider for usages
func (k *Kernel) initProvider(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not retrieve working directory, error: %s", err.Error())
	}

	log.Println("Loading provider plugin configuration...")
	pluginCfgErr := k.p.LoadConfig(wd + "/config.json")
	if pluginCfgErr != nil {
		return fmt.Errorf("provider configuration, error: %s", pluginCfgErr.Error())
	}

	log.Println("Preparing provider plugin...")
	pluginBootErr := k.p.Boot()
	if pluginBootErr != nil {
		return fmt.Errorf("provider boot, error:: %s", pluginBootErr.Error())
	}

	return nil
}

// initDashboard to display monitored system data
func (k *Kernel) initDashboard(ctx context.Context, t terminalapi.Terminal) error {
	// TODO Can be added provider name to dashboard, interface update required
	log.Println("Fetching data graph from provider...")
	g, gErr := k.p.DispatchGraph()
	if gErr != nil {
		return fmt.Errorf("provider dispatch monitoring data failed, error: %s", gErr.Error())
	}

	log.Println("Creating terminal dashboard UI...")
	var dErr error
	k.d, dErr = dashboard.NewMonitoringDashboard("BitAccretion", t, g)
	if dErr != nil {
		return dErr
	}

	k.l = k.d.EventLogger
	k.o = event.NewDashboardObserver(k.d.EventLogger)
	k.o.RegisterNewSubscriber(k.d)

	return nil
}

// dashboardUpdate ask provider to collect new data and push update to widgets
func (k *Kernel) dashboardUpdate(ctx context.Context, delay time.Duration) {
	isNeedToFetch := true
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if isNeedToFetch == false {
				continue
			}

			isNeedToFetch = false
			k.l.Normal("Requesting provider to get new data update...")

			providerGraph, providerGraphErr := k.p.FetchNewData()
			if providerGraphErr != nil {
				k.l.Error(providerGraphErr.Error())
				return
			}

			k.o.NotifySubscribers(event.UpdateEvent{MonitoringGraph: providerGraph})
			isNeedToFetch = true
		case <-ctx.Done():
			return
		}
	}
}

// Run main process to handle dashboard and update it with data from provider
func (k *Kernel) Run(t terminalapi.Terminal) error {
	log.Println("Initializing kernel...")

	ctx, cancel := context.WithCancel(context.Background())
	providerErr := k.initProvider(ctx)
	if providerErr != nil {
		return providerErr
	}

	dashErr := k.initDashboard(ctx, t)
	if dashErr != nil {
		return dashErr
	}
	log.Println("Kernel ready...")
	log.Println("Rendering terminal UI...")
	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	// TODO Time update must be used from Kernel cfg
	go k.dashboardUpdate(ctx, 1*time.Second)

	fmt.Print("\033[H\033[2J") // Clean terminal screen from any artifacts

	// TODO Time terminal UI redraw must be used from Kernel cfg
	termErr := termdash.Run(
		ctx,
		t,
		k.d.TerminalContainer,
		termdash.KeyboardSubscriber(quitter),
		termdash.RedrawInterval(1*time.Second),
	)
	if termErr != nil {
		return termErr
	}

	return nil
}
