package event

import (
	"github.com/LinMAD/BitAccretion/logger"
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

func (s *testObserverSubscriber) HandleNotifyEvent(e UpdateEvent) error {
	if s.GetName() == testSubscriberName {
		testSubscriberValid = true
	}

	testSubscriberNotified = true

	return nil
}

func (s *testObserverSubscriber) GetName() string {
	return s.name
}

func TestObserverSubscriber(t *testing.T) {
	ob := NewDashboardObserver(logger.NullLogger{})
	ob.RegisterNewSubscriber(&testObserverSubscriber{name: testSubscriberName})

	ob.NotifySubscribers(UpdateEvent{})
	if !testSubscriberNotified || !testSubscriberValid {
		t.FailNow()
	}
}
