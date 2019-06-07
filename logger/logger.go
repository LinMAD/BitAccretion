package logger

// LevelOfLog message leveling
type LevelOfLog int

const (
	// DebugLog allows send in debug channel
	DebugLog LevelOfLog = iota
	// NormalLog allows write to errors and normal log but debug will be excluded
	NormalLog
)

// ILogger basic interface
type ILogger interface {
	// Mode of dashboardLogger
	SetMode(level LevelOfLog)
	// Debug messages
	Debug(msg string)
	// Normal general messages
	Normal(msg string)
	// Error messages
	Error(msg string)
}

// NullLogger can be used for tests, not writing data
type NullLogger struct{}

// SetMode of logger
func (n NullLogger) SetMode(level LevelOfLog) {
	return
}

// Debug chanel
func (n NullLogger) Debug(msg string) {
	return
}

// Normal chanel
func (n NullLogger) Normal(msg string) {
	return
}

// Error chanel
func (n NullLogger) Error(msg string) {
	return
}
