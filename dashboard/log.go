package dashboard

import (
	"github.com/LinMAD/BitAccretion/logger"
	"github.com/mum4k/termdash/cell"
)

// loggerHandler used to write log events
type loggerHandler struct {
	lvl    logger.LevelOfLog
	widget *TextWidgetHandler
}

// SetTextWidget where log messages delivered
func (l *loggerHandler) SetTextWidget(textScreen *TextWidgetHandler) {
	l.widget = textScreen
}

// SetMode of loggerHandler
func (l *loggerHandler) SetMode(level logger.LevelOfLog) {
	l.lvl = level
}

// Debug messages
func (l *loggerHandler) Debug(msg string) {
	if l.lvl == logger.DebugLog {
		l.widget.WriteToEventLog(msg, cell.ColorWhite)
	}
}

// Normal events messages
func (l *loggerHandler) Normal(msg string) {
	l.widget.WriteToEventLog(msg, cell.ColorWhite)
}

// Error events messages
func (l *loggerHandler) Error(msg string) {
	l.widget.WriteToEventLog(msg, cell.ColorRed)
}
