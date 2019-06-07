package model

// HealthState flag
type HealthState int

const (
	// HealthNormal ok flag
	HealthNormal HealthState = iota
	// HealthWarning some errors
	HealthWarning
	// HealthCritical there is issue
	HealthCritical
)

// HealthStatesMap to string
var HealthStatesMap = map[HealthState]string{
	HealthNormal:   "Normal",
	HealthWarning:  "Warning",
	HealthCritical: "Critical",
}

// HealthSensitivity represents metrics sensitivity determination
type HealthSensitivity struct {
	Danger  int `json:"danger"`
	Warning int `json:"warning"`
}
