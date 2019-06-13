package util

import "github.com/LinMAD/BitAccretion/model"

// GetMetricHealthByValue return node health by error limits
func GetMetricHealthByValue(m *model.SystemMetric, h *model.HealthSensitivity) model.HealthState {
	return dispatchHealth(int(m.ErrorCount), h)
}

// GetMetricsHealthByPercentRatio provide node health by percent ration with ok req and errors
func GetMetricsHealthByPercentRatio(m *model.SystemMetric, h *model.HealthSensitivity) model.HealthState {
	return dispatchHealth(int((m.ErrorCount/m.RequestCount)*100), h)
}

// dispatchHealth by ratio and HealthSensitivity
func dispatchHealth(ratio int, h *model.HealthSensitivity) model.HealthState {
	if ratio >= h.Critical {
		return model.HealthCritical
	} else if ratio >= h.Warning {
		return model.HealthWarning
	} else {
		return model.HealthNormal
	}
}
