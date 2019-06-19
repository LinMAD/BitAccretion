package model

import "github.com/LinMAD/BitAccretion/logger"

// Config core settings
type Config struct {
	// SoundAlertDelayMin between playing sound
	SoundAlertDelayMin int `json:"sound_alert_delay_min"`
	// IsSoundMode enabled for alerts
	IsSoundMode bool `json:"sound_mode"`
	// SurveyIntervalSec for data updates
	SurveyIntervalSec int `json:"survey_interval_sec"`
	// InterfaceUpdateIntervalSec terminal redraw frequency
	InterfaceUpdateIntervalSec int `json:"interface_update_interval_sec"`
	// LogLevel message to display in event widget
	LogLevel logger.LevelOfLog `json:"log_level"`
	// DisplayEvenLogHistory limit to display rendered charts or logs
	DisplayEvenLogHistory int16 `json:"display_even_log_history"`
}
