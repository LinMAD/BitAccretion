package dashboard

import (
	"fmt"
	"time"

	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/extension"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

const maxTextHistory = 100

// AnnouncerHandler for dashboard
type AnnouncerHandler struct {
	name           string
	t              *text.Text
	s              extension.ISound
	historyCounter int8
	lastSoundTriggerTime time.Time
}

// HandleNotifyEvent write to text widget
func (anon *AnnouncerHandler) HandleNotifyEvent(e event.UpdateEvent) error {
	healthMsgList := make(map[model.HealthState]string, 0)
	systems := e.MonitoringGraph.GetAllVertices()

	anon.handleHistory()
	for i := 0; i < len(systems); i++ {
		sys := systems[i]

		if sys.Health != model.HealthNormal {
			healthMsgList[sys.Health] = fmt.Sprintf("%s - Health: %s", sys.Name, model.HealthStatesMap[sys.Health])
		}

		if sys.Health == model.HealthCritical {
			anon.playAlter(sys.Name)
		}
	}

	for h, msg := range healthMsgList {
		var termColor cell.Color

		switch h {
		case model.HealthCritical:
			termColor = cell.ColorRed
		case model.HealthWarning:
			termColor = cell.ColorYellow
		default:
			termColor = cell.ColorWhite
		}

		anon.WriteToEventLog(msg, termColor)
	}

	return nil
}

// GetName of widget handler
func (anon *AnnouncerHandler) GetName() string {
	return anon.name
}

// WriteToEventLog display message with color in widget
func (anon *AnnouncerHandler) WriteToEventLog(msg string, color cell.Color) {
	writeErr := anon.t.Write(fmt.Sprintf("|%s| %s\n", time.Now().Format(time.Stamp), msg), text.WriteCellOpts(cell.FgColor(color)))
	if writeErr != nil {
		panic(writeErr)
	}
}

// handleHistory of logged messages
func (anon *AnnouncerHandler) handleHistory() {
	anon.historyCounter++
	if anon.historyCounter <= maxTextHistory {
		return
	}

	anon.historyCounter = 0
	anon.t.Reset()
}

// playAlter sound for given name
func (anon *AnnouncerHandler) playAlter(name string) {
	if anon.s == nil {
		return
	}

	now := time.Now().UTC()
	_ = now.Sub(anon.lastSoundTriggerTime)

	anon.WriteToEventLog(fmt.Sprintf("Playing alert sound for %s...", name), cell.ColorBlue)
	anon.s.PlayAlert(model.VertexName(name))
	anon.lastSoundTriggerTime = now
}

// NewAnnouncerWidget creates and returns prepared widget
func NewAnnouncerWidget(sound extension.ISound, name string) (*AnnouncerHandler, error) {
	t, tErr := text.New(text.WrapAtRunes(), text.WrapAtWords(), text.RollContent())
	if tErr != nil {
		return nil, tErr
	}

	return &AnnouncerHandler{name: name, t: t, s: sound, lastSoundTriggerTime: time.Now().UTC()}, nil
}
