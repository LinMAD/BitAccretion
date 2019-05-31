package event

import (
	"github.com/LinMAD/BitAccretion/model"
)

// TODO Fix interface naming (they not attached to widgets any more)

type (
	UpdateEvent struct {
		// MonitoringGraph represents a structure for dashboard
		MonitoringGraph model.Graph
	}

	// IWidgetSubscriber represent basic interface to update dashboard on event
	IWidgetSubscriber interface {
		// HandleNotifyEvent allows to publish update in subscriber
		HandleNotifyEvent(UpdateEvent)
		// GetName returns name of subscriber
		GetName() string
	}

	// IWidgetObserver represents implementation to update dashboard widgets
	IWidgetObserver interface {
		// RegisterSubscriber widget observer
		RegisterSubscriber(IWidgetSubscriber)
		// NotifySubscribers publishes new events to listeners\subscribers
		NotifySubscribers(UpdateEvent)
	}

	// widgetObserver must be used to deliver updates to different subscribed widgets
	widgetObserver struct {
		widgetSubscribers []IWidgetSubscriber
	}
)

// NewDashboardObserver returns new instance of observer for widgets
func NewDashboardObserver() IWidgetObserver {
	return &widgetObserver{
		widgetSubscribers: make([]IWidgetSubscriber, 0),
	}
}

// RegisterSubscriber to observer
func (wn *widgetObserver) RegisterSubscriber(wo IWidgetSubscriber) {
	wn.widgetSubscribers = append(wn.widgetSubscribers, wo)
}

// NotifySubscribers all subscribers
func (wn *widgetObserver) NotifySubscribers(event UpdateEvent) {
	for _, o := range wn.widgetSubscribers {
		o.HandleNotifyEvent(event)
	}
}
