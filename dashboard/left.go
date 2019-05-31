package dashboard

import (
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/mum4k/termdash/widgets/barchart"
	"math/rand"
	"time"
)

func CreateLeftLayout(systems []model.Node) (container.LeftOption, error) {
	sysBar, sysBarErr := createSystemBar(systems)
	if sysBarErr != nil {
		return nil, sysBarErr
	}

	leftLayout := container.Left(
		container.Border(linestyle.Round),
		container.BorderTitle("Requests to systems"),
		container.PlaceWidget(sysBar),
	)

	return leftLayout, nil
}

func createSystemBar(systems []model.Node) (widgetapi.Widget, error) {
	barWidth := 0
	sysCount := len(systems)
	sysNames := make([]string, sysCount)
	sysBarColors := make([]cell.Color, sysCount)
	sysValBarColors := make([]cell.Color, sysCount)

	for i := 0; i < sysCount; i++ {
		sysNames[i] = systems[i].Name

		switch systems[i].Health {
		case model.Warning:
			sysBarColors[i] = cell.ColorYellow
			sysValBarColors[i] = cell.ColorBlack
		case model.Critical:
			sysBarColors[i] = cell.ColorRed
			sysValBarColors[i] = cell.ColorBlack
		default:
			sysBarColors[i] = cell.ColorBlue
			sysValBarColors[i] = cell.ColorBlack
		}

		if barWidth < len(systems[i].Name) {
			barWidth = len(systems[i].Name)
		}
	}

	sysBar, sysBarErr := barchart.New(
		barchart.BarColors(sysBarColors),
		barchart.ValueColors(sysValBarColors),
		barchart.ShowValues(),
		barchart.Labels(sysNames),
	)
	if sysBarErr != nil {
		return nil, sysBarErr
	}

	go playBarChart(sysBar, systems, 1*time.Second)

	return sysBar, nil
}

func playBarChart(bc *barchart.BarChart, systems []model.Node, delay time.Duration) {
	sysMaxValue := 0

	for i := 0; i < len(systems); i++ {
		if sysMaxValue < systems[i].Metric.RequestCount {
			sysMaxValue = systems[i].Metric.RequestCount
		}
	}

	ticker := time.NewTicker(delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var values []int
			for i := 0; i < len(systems); i++ {
				values = append(values, int(rand.Int31n(int32(sysMaxValue))))
			}

			if err := bc.Values(values, sysMaxValue); err != nil {
				panic(err)
			}
		}
	}
}
