package model

// SystemMetric basic metric data of observable system
type SystemMetric struct {
	// RequestCount of system
	RequestCount int
	// ErrorCount of system
	ErrorCount int
}

// SystemMetaData additional information in string format
type SystemMetaData interface{}
