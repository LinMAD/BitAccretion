package model

// HealthState flag
type HealthState int

const (
	Normal HealthState = iota
	Warning
	Critical
)

// HealthStatesMap to string
var HealthStatesMap = map[HealthState]string{
	Normal: "Normal",
	Warning: "Warning",
	Critical: "Critical",
}
