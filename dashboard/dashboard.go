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

// MonitoringDashboard core dashboard structure with constructed widgets
type MonitoringDashboard struct {
	observer          event.IObserver
	TerminalContainer *container.Container
}

// HandleNotifyEvent send update to monitoring dashboard
func (m MonitoringDashboard) HandleNotifyEvent(e event.UpdateEvent) {
	m.observer.NotifySubscribers(e)
}

// GetName of subscriber
func (m MonitoringDashboard) GetName() string {
	return "MonitoringDashboard"
}

// NewMonitoringDashboard with constructed widgets
func NewMonitoringDashboard(dashboardName string, t terminalapi.Terminal) (*MonitoringDashboard, error) {
	// TODO Split method to widgets init

	// TODO Add factory of charts
	// Init widgets
	okReqsBarWidget, okReqsBarWidgetErr := NewBarChart("ok_reqs_bar", cell.ColorBlue, true, stub.GetStubNodes())
	if okReqsBarWidgetErr != nil {
		return nil, okReqsBarWidgetErr
	}

	badReqsBarWidget, badReqsBarWidgetErr := NewBarChart("bad_reqs_bar", cell.ColorRed, false, stub.GetStubNodes())
	if badReqsBarWidgetErr != nil {
		return nil, okReqsBarWidgetErr
	}

	aggSparkSuccessReq, aggSparkSuccessReqErr := NewSparkLineChart("aggregated_reqs_in_line")
	if aggSparkSuccessReqErr != nil {
		return nil, aggSparkSuccessReqErr
	}

	// Create dashboard observer
	termDash := &MonitoringDashboard{
		observer: event.NewDashboardObserver(),
	}

	// Register widgets
	termDash.observer.RegisterNewSubscriber(okReqsBarWidget)
	termDash.observer.RegisterNewSubscriber(badReqsBarWidget)
	termDash.observer.RegisterNewSubscriber(aggSparkSuccessReq)

	// TODO Spilt method to layout construct (left and right)

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle(dashboardName),
		container.SplitHorizontal(
			container.Top(
				container.SplitVertical(
					container.Left(
						container.Border(linestyle.Light),
						stubEventLogWidget(),
					),
					container.Right(
						container.Border(linestyle.Light),
						container.BorderTitle("Aggregated requests"),
						container.PlaceWidget(aggSparkSuccessReq.lc),
					),
				),
			),
			container.Bottom(
				container.Border(linestyle.Round),
				container.BorderTitle("Requests to systems"),
				container.SplitVertical(
					container.Left(
						container.Border(linestyle.Light),
						container.BorderTitle("Successful"),
						container.PlaceWidget(okReqsBarWidget.barChart),
					),
					container.Right(
						container.Border(linestyle.Light),
						container.BorderTitle("Incorrect"),
						container.PlaceWidget(badReqsBarWidget.barChart),
					),
				),
			),
		),
	)

	termDash.TerminalContainer = c

	return termDash, err
}

// TODO Remove that stub
func stubEventLogWidget() container.Option {
	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:50| Error in Name 1\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:51| Error in Name 1\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:52| Error in Name 1\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:53| Error in Name 2\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}

	return container.PlaceWidget(wrapped)
}
