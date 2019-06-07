package dashboard

import (
	"strings"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
)

// ClockWidgetHandler for dashboard
type ClockWidgetHandler struct {
	sdClock *segmentdisplay.SegmentDisplay
}

// runClock execute sdClock update
func (c *ClockWidgetHandler) runClock() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			nowStr := now.Format("15:04:05")
			parts := strings.Split(nowStr, ":")

			spacer := " "
			if now.Second()%2 == 0 {
				spacer = ":"
			}

			chunks := []*segmentdisplay.TextChunk{
				segmentdisplay.NewChunk(parts[0], segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorWhite))),
				segmentdisplay.NewChunk(spacer, segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorYellow))),
				segmentdisplay.NewChunk(parts[1], segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorWhite))),
				segmentdisplay.NewChunk(spacer, segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorYellow))),
				segmentdisplay.NewChunk(parts[2], segmentdisplay.WriteCellOpts(cell.FgColor(cell.ColorWhite))),
			}

			if err := c.sdClock.Write(chunks); err != nil {
				panic(err)
			}
		}
	}
}

// NewClockWidget creates and returns prepared widget
func NewClockWidget() (*ClockWidgetHandler, error) {
	c, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}

	cwh := &ClockWidgetHandler{sdClock: c}
	go cwh.runClock()

	return cwh, nil
}
