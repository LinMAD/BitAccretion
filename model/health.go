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
	HealthNormal:   "HealthNormal",
	HealthWarning:  "HealthWarning",
	HealthCritical: "HealthCritical",
}
