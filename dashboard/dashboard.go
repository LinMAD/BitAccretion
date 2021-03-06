package dashboard

import (
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/extension"
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
	regression    *GaugeRegressHandler
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
func (m *MonitoringDashboard) initWidgets(s extension.ISound, c *model.Config, n []*model.Node) (err error) {
	m.widgetCollection.reqSuccessful, err = NewBarWidget("ok_reqs_bar_widget", cell.ColorGreen, true, n)
	if err != nil {
		return err
	}

	m.widgetCollection.reqIncorrect, err = NewBarWidget("bad_reqs_bar_widget", cell.ColorRed, false, n)
	if err != nil {
		return err
	}

	m.widgetCollection.reqAggregated, err = NewLineWidget("aggregated_reqs_line_widget", c)
	if err != nil {
		return err
	}

	m.widgetCollection.eventLog, err = NewAnnouncerWidget(s, c, "system_error_text_widget")
	if err != nil {
		return err
	}

	m.widgetCollection.regression, err = NewRegressionWidget("Regression level of errors", c)
	if err != nil {
		return err
	}

	return nil
}

// createLayout for dashboard and place widgets
func (m *MonitoringDashboard) createLayout(dashboardName string, t *terminalapi.Terminal) (err error) {
	m.TerminalContainer, err = container.New(
		*t,
		container.Border(linestyle.Light),
		container.BorderTitle(dashboardName),
		container.SplitHorizontal(
			container.Top(
				container.SplitVertical(
					container.Left(
						container.SplitHorizontal(
							container.Top(
								container.Border(linestyle.Round),
								container.PlaceWidget(m.widgetCollection.regression.gauge),
							),
							container.Bottom(
								container.Border(linestyle.Round),
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
func NewMonitoringDashboard(dashboardName string, c *model.Config, s extension.ISound, t terminalapi.Terminal, g model.Graph) (*MonitoringDashboard, error) {
	termDash := &MonitoringDashboard{widgetCollection: &widgets{}}

	initErr := termDash.initWidgets(s, c, g.GetAllVertices())
	if initErr != nil {
		return nil, initErr
	}

	layoutErr := termDash.createLayout(dashboardName, &t)
	if layoutErr != nil {
		return nil, layoutErr
	}

	// Add dependencies
	termDash.EventLogger = &loggerHandler{lvl: c.LogLevel, widget: termDash.widgetCollection.eventLog}
	termDash.observer = event.NewDashboardObserver(termDash.EventLogger)

	// Register widgets to be observable for g updates
	termDash.observer.RegisterNewSubscriber(termDash.widgetCollection.reqSuccessful)
	termDash.observer.RegisterNewSubscriber(termDash.widgetCollection.reqIncorrect)
	termDash.observer.RegisterNewSubscriber(termDash.widgetCollection.reqAggregated)
	termDash.observer.RegisterNewSubscriber(termDash.widgetCollection.eventLog)
	termDash.observer.RegisterNewSubscriber(termDash.widgetCollection.regression)

	return termDash, nil
}
