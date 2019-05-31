package dashboard

import (
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/stub"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/text"
)

// MonitoringDashboard
type MonitoringDashboard struct {
	observer          event.IObserver
	TerminalContainer *container.Container
}

// HandleNotifyEvent send update to monitoring dashboard
func (m MonitoringDashboard) HandleNotifyEvent(e event.UpdateEvent) {
	m.observer.NotifySubscribers(e)
}

// GetName
func (m MonitoringDashboard) GetName() string {
	return "MonitoringDashboard"
}

// NewMonitoringDashboard with constructed widgets
func NewMonitoringDashboard(t terminalapi.Terminal) (*MonitoringDashboard, error) {
	// TODO Split method to widgets init

	// Init widgets
	barWidget, barWidgetErr := NewBarChart("BarChartWidget", stub.GetStubNodes())
	if barWidgetErr != nil {
		return nil, barWidgetErr
	}

	// Create dashboard observer
	termDash := &MonitoringDashboard{
		observer: event.NewDashboardObserver(),
	}

	// Register widgets
	termDash.observer.RegisterNewSubscriber(barWidget)

	// TODO Spilt method to layout construct (left and right)

	leftLayout := container.Left(
		container.Border(linestyle.Round),
		container.BorderTitle("Requests to systems"),
		container.PlaceWidget(barWidget.barChart),
	)

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.SplitVertical(
			leftLayout,
			container.Right(
				container.SplitHorizontal(
					container.Top(
						container.Border(linestyle.Light),
						EventLogWidget(),
					),
					container.Bottom(
						container.Border(linestyle.Light),
					),
				),
			),
		),
	)

	termDash.TerminalContainer = c

	return termDash, err
}

func EventLogWidget() container.Option {
	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:50| Error in CSEye\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:51| Error in CSEye\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:52| Error in CSEye\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:53| Error in Fastlane\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}

	return container.PlaceWidget(wrapped)
}
