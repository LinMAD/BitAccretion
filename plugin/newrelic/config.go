package newrelic

// NRConfig addition for config required for New Relic
type NRConfig struct {
	APIKey            string                     `json:"api_key"`
	APPSets           []APPSet                   `json:"app_sets"`
	EnabledSoundAlert bool                       `json:"enabled_sound_alert"`
	SurveyTime        int                        `json:"survey_time"`
}

// APPSet registered in New Relic
type APPSet struct {
	AppDetails
	Nested []APPSet `json:"nested"`
}

// AppDetails detailed information about application in New Relic
type AppDetails struct {
	Name         string   `json:"name"`
	ID           string   `json:"id"`
	RelicMetrics []string `json:"relic_metrics"`
}
