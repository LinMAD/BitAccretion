package dashboard

import (
	"fmt"
	"github.com/LinMAD/BitAccretion/event"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
	"time"
)

const maxTextHistory = 100

// TextWidgetHandler for dashboard
type TextWidgetHandler struct {
	name           string
	t              *text.Text
	historyCounter int8
}

func (txt *TextWidgetHandler) HandleNotifyEvent(e event.UpdateEvent) {
	// TODO Use one more function from TextWidgetHandler struct to have different event writers not only for graph

	healthMsgList := make(map[model.HealthState]string, 0)
	systems := e.MonitoringGraph.GetAllVertices()
	txt.historyCounter++

	for i := 0; i < len(systems); i++ {
		sys := systems[i]

		if sys.Health != model.HealthNormal {
			healthMsgList[sys.Health] = fmt.Sprintf(
				"|%s| %s - Health: %s\n",
				time.Now().Format(time.Stamp),
				sys.Name,
				model.HealthStatesMap[sys.Health],
			)
		}
	}

	if txt.historyCounter >= maxTextHistory {
		txt.historyCounter = 0
		txt.t.Reset()
	}

	for h, msg := range healthMsgList {
		var writeErr error

		if h == model.HealthCritical {
			writeErr = txt.t.Write(msg, text.WriteCellOpts(cell.FgColor(cell.ColorRed)))
		} else if h == model.HealthWarning {
			writeErr = txt.t.Write(msg, text.WriteCellOpts(cell.FgColor(cell.ColorYellow)))
		}

		if writeErr != nil {
			panic(writeErr) // TODO Handle in grace way, log or ignore
		}
	}
}

// GetName of widget handler
func (txt *TextWidgetHandler) GetName() string {
	return txt.name
}

// NewTextWidget creates and returns prepared widget
func NewTextWidget(name string) (*TextWidgetHandler, error) {
	t, tErr := text.New(text.WrapAtRunes(), text.WrapAtWords(), text.RollContent())
	if tErr != nil {
		return nil, tErr
	}

	return &TextWidgetHandler{name: name, t: t}, nil
}
