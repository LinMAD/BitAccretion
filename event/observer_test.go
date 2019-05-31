package event

import (
	"testing"
)

const testSubscriberName = "test_sub"

var (
	testSubscriberNotified = false
	testSubscriberValid    = false
)

type testObserverSubscriber struct {
	name string
}

func (s *testObserverSubscriber) HandleNotifyEvent(e UpdateEvent) {
	if s.GetName() == testSubscriberName {
		testSubscriberValid = true
	}

	testSubscriberNotified = true
}

func (s *testObserverSubscriber) GetName() string {
	return s.name
}

func TestObserverSubscriber(t *testing.T) {
	ob := NewDashboardObserver()
	RegisterSubscriber(&testObserverSubscriber{
		name: testSubscriberName,
	})

	NotifySubscribers(UpdateEvent{})
	if !testSubscriberNotified || !testSubscriberValid {
		t.FailNow()
	}
}
