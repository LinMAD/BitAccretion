package model

// SystemMetric basic metric data of observable system
type SystemMetric struct {
	// RequestCount of system
	RequestCount float32
	// ErrorCount of system
	ErrorCount float32
}

// SystemMetaData additional information in string format
type SystemMetaData interface{}
