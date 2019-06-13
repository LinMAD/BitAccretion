package util

import "github.com/LinMAD/BitAccretion/model"

// GetMetricHealthByValue return node health by error limits
func GetMetricHealthByValue(m *model.SystemMetric, sense *model.HealthSensitivity) model.HealthState {
	return dispatchHealth(int(m.ErrorCount), sense)
}

// GetMetricsHealthByPercentRatio provide node health by percent ration with ok req and errors
func GetMetricsHealthByPercentRatio(m *model.SystemMetric, sense *model.HealthSensitivity) model.HealthState {
	return dispatchHealth(int((m.ErrorCount / m.RequestCount) * 100), sense)
}

// dispatchHealth by ratio and HealthSensitivity
func dispatchHealth(ratio int, sense *model.HealthSensitivity) model.HealthState {
	if ratio >= sense.Danger {
		return model.HealthCritical
	} else if ratio >= sense.Warning {
		return model.HealthWarning
	} else {
		return model.HealthNormal
	}
}
