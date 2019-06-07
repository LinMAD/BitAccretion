package model

// Config core settings
type Config struct {
	// SurveyIntervalSec for data updates
	SurveyIntervalSec int `json:"survey_interval_sec"`
	// InterfaceUpdateIntervalSec terminal redraw frequency
	InterfaceUpdateIntervalSec int `json:"interface_update_interval_sec"`
}
