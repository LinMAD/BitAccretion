package dashboard

import (
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/logger"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

// MonitoringDashboard core dashboard structure with constructed widgets
type MonitoringDashboard struct {
	TerminalContainer *container.Container
	EventLogger       logger.ILogger
	observer          event.IObserver
	widgetCollection  *widgets
}

// widgets of dashboard
type widgets struct {
	reqSuccessful *BarWidgetHandler
	reqIncorrect  *BarWidgetHandler
	reqAggregated *SparkLineWidgetHandler
	eventLog      *AnnouncerHandler
	clock         *ClockWidgetHandler
}

// HandleNotifyEvent send update to monitoring dashboard
func (m *MonitoringDashboard) HandleNotifyEvent(e event.UpdateEvent) error {
	m.observer.NotifySubscribers(e)

	return nil
}

// GetName of subscriber
func (m *MonitoringDashboard) GetName() string {
	return "monitoring_dashboard"
}

// initWidgets for dashboard
func (m *MonitoringDashboard) initWidgets(nodes []*model.Node) (err error) {
	m.widgetCollection.reqSuccessful, err = NewBarWidget("ok_reqs_bar_widget", cell.ColorGreen, true, nodes)
	if err != nil {
		return err
	}

	m.widgetCollection.reqIncorrect, err = NewBarWidget("bad_reqs_bar_widget", cell.ColorRed, false, nodes)
	if err != nil {
		return err
	}

	m.widgetCollection.reqAggregated, err = NewLineWidget("aggregated_reqs_line_widget")
	if err != nil {
		return err
	}

	m.widgetCollection.eventLog, err = NewAnnouncerWidget("system_error_text_widget")
	if err != nil {
		return err
	}
	m.widgetCollection.clock, err = NewClockWidget()
	if err != nil {
		return err
	}

	return nil
}

// createLayout for dashboard and place widgets
func (m *MonitoringDashboard) createLayout(dashboardName string, t *terminalapi.Terminal) (err error) {
	m.TerminalContainer, err = container.New(
		*t,
		container.Border(linestyle.Double),
		container.BorderTitle(dashboardName),
		container.SplitHorizontal(
			container.Top(
				container.SplitVertical(
					container.Left(
						container.SplitHorizontal(
							container.Top(
								container.Border(linestyle.Round),
								container.PlaceWidget(m.widgetCollection.clock.sdClock),
							),
							container.Bottom(
								container.Border(linestyle.Double),
								container.BorderTitle("Event log"),
								container.PlaceWidget(m.widgetCollection.eventLog.t),
							),
						),
					),
					container.Right(
						container.Border(linestyle.Light),
						container.BorderTitle("Aggregated requests"),
						container.PlaceWidget(m.widgetCollection.reqAggregated.lc),
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
						container.PlaceWidget(m.widgetCollection.reqSuccessful.barChart),
					),
					container.Right(
						container.Border(linestyle.Light),
						container.BorderTitle("Incorrect"),
						container.PlaceWidget(m.widgetCollection.reqIncorrect.barChart),
					),
				),
			),
		),
	)

	return err
}

// NewMonitoringDashboard constructor, will prepare widgets, subscriber's and dependencies
func NewMonitoringDashboard(dashboardName string, logLvl logger.LevelOfLog, t terminalapi.Terminal, graph model.Graph) (*MonitoringDashboard, error) {
	termDash := &MonitoringDashboard{widgetCollection: &widgets{}}

	initErr := termDash.initWidgets(graph.GetAllVertices())
	if initErr != nil {
		return nil, initErr
	}

	layoutErr := termDash.createLayout(dashboardName, &t)
	if layoutErr != nil {
		return nil, layoutErr
	}

	// Add dependencies
	termDash.EventLogger = &loggerHandler{lvl: logLvl, widget: termDash.widgetCollection.eventLog}
	termDash.observer = event.NewDashboardObserver(termDash.EventLogger)

	// Register widgets to be observable for graph updates
	termDash.observer.RegisterNewSubscriber(termDash.widgetCollection.reqSuccessful)
	termDash.observer.RegisterNewSubscriber(termDash.widgetCollection.reqIncorrect)
	termDash.observer.RegisterNewSubscriber(termDash.widgetCollection.reqAggregated)
	termDash.observer.RegisterNewSubscriber(termDash.widgetCollection.eventLog)

	return termDash, nil
}
