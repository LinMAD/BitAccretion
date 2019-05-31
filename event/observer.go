package event

import (
	"github.com/LinMAD/BitAccretion/model"
)

type (
	UpdateEvent struct {
		// MonitoringGraph represents a structure for dashboard
		MonitoringGraph model.Graph
	}

	// ISubscriber represent basic interface to update dashboard on event
	ISubscriber interface {
		// HandleNotifyEvent allows to publish update in subscriber
		HandleNotifyEvent(UpdateEvent)
		// GetName returns name of subscriber
		GetName() string
	}

	// IObserver represents implementation to update dashboard widgets
	IObserver interface {
		// RegisterNewSubscriber widget observer
		RegisterNewSubscriber(ISubscriber)
		// NotifySubscribers publishes new events to listeners\subscribers
		NotifySubscribers(UpdateEvent)
	}

	// observer must be used to deliver updates to different subscribed widgets
	observer struct {
		subscribers []ISubscriber
	}
)

// NewDashboardObserver returns new instance of observer for widgets
func NewDashboardObserver() IObserver {
	return &observer{
		subscribers: make([]ISubscriber, 0),
	}
}

// RegisterNewSubscriber to observer
func (wn *observer) RegisterNewSubscriber(wo ISubscriber) {
	wn.subscribers = append(wn.subscribers, wo)
}

// NotifySubscribers all subscribers
func (wn *observer) NotifySubscribers(event UpdateEvent) {
	for _, o := range wn.subscribers {
		o.HandleNotifyEvent(event)
	}
}
