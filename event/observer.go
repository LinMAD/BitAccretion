package event

import (
	"fmt"
	"github.com/LinMAD/BitAccretion/logger"
	"github.com/LinMAD/BitAccretion/model"
)

type (
	// UpdateEvent general event structure of subscribers
	UpdateEvent struct {
		// MonitoringGraph represents a structure for dashboard
		MonitoringGraph model.Graph
	}

	// ISubscriber represent basic interface to update dashboard on event
	ISubscriber interface {
		// HandleNotifyEvent allows to publish update in subscriber
		HandleNotifyEvent(UpdateEvent) error
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
		log         logger.ILogger
		subscribers []ISubscriber
	}
)

// NewDashboardObserver returns new instance of observer for widgets
func NewDashboardObserver(logger logger.ILogger) IObserver {
	return &observer{log: logger, subscribers: make([]ISubscriber, 0)}
}

// RegisterNewSubscriber to observer
func (o *observer) RegisterNewSubscriber(wo ISubscriber) {
	o.subscribers = append(o.subscribers, wo)
}

// NotifySubscribers all subscribers
func (o *observer) NotifySubscribers(event UpdateEvent) {
	for _, s := range o.subscribers {
		err := s.HandleNotifyEvent(event)
		if err != nil {
			o.log.Error(fmt.Sprintf("%s has error -> %s", s.GetName(), err.Error()))
		}
	}
}
