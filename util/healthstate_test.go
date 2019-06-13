package util

import (
	"github.com/LinMAD/BitAccretion/model"
	"github.com/stretchr/testify/suite"
	"testing"
)

type HealthStateTestSuite struct {
	suite.Suite
	sense *model.HealthSensitivity
}

func (t *HealthStateTestSuite) SetupTest() {
	t.sense = &model.HealthSensitivity{
		Critical: 20,
		Warning:  10,
	}
}

func TestRunGraphTestSuite(t *testing.T) {
	suite.Run(t, new(HealthStateTestSuite))
}

func (t *HealthStateTestSuite) TestGetMetricHealthByValue() {
	var s model.HealthState

	s = GetMetricHealthByValue(&model.SystemMetric{RequestCount: 100, ErrorCount: 50}, t.sense)
	if s != model.HealthCritical {
		t.Fail("Expected to be 'Critical' state if errors 50 but limit 20")
	}

	s = GetMetricHealthByValue(&model.SystemMetric{RequestCount: 100, ErrorCount: 10}, t.sense)
	if s != model.HealthWarning {
		t.Fail("Expected to be 'Warning' state if errors 10 but limit 10")
	}

	s = GetMetricHealthByValue(&model.SystemMetric{RequestCount: 100, ErrorCount: 5}, t.sense)
	if s != model.HealthNormal {
		t.Fail("Expected to be 'Normal' state if errors 5 but limit not reached")
	}
}

func (t *HealthStateTestSuite) TestGetMetricsHealthByPercentRatio() {
	var s model.HealthState

	s = GetMetricHealthByValue(&model.SystemMetric{RequestCount: 100, ErrorCount: 50}, t.sense)
	if s != model.HealthCritical {
		t.Fail("Expected to be 'Critical' state if 50% = (50/100) * 100, limit 20%")
	}

	s = GetMetricHealthByValue(&model.SystemMetric{RequestCount: 100, ErrorCount: 10}, t.sense)
	if s != model.HealthWarning {
		t.Fail("Expected to be 'Critical' state if 10% = (10/100) * 100, limit 10%")
	}

	s = GetMetricHealthByValue(&model.SystemMetric{RequestCount: 100, ErrorCount: 5}, t.sense)
	if s != model.HealthNormal {
		t.Fail("Expected to be 'Critical' state if 5% = (5/100) * 100, limit not reached")
	}
}
