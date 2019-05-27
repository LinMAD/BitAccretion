package dashboard

import (
	"fmt"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/mum4k/termdash/widgets/linechart"
	"math"
)

func CreateLeftLayout(systems []model.Node) (container.LeftOption, error) {
	var leftGrid container.LeftOption
	sysCount := len(systems)

	widgetsList := make([][]container.Option, 0)

	for i := 0; i < sysCount; i++ {
		widgetOptions := make([]container.Option, 0)
		w, wErr := generateWidgets(systems[i])
		if wErr != nil {
			return nil, wErr
		}

		widgetOptions = append(widgetOptions, w...)
		widgetsList = append(widgetsList, widgetOptions)
	}

	leftGrid = container.Left(
		container.SplitHorizontal(
			container.Top(widgetsList[0]...),
			container.Bottom(widgetsList[1]...),
		),
	)

	return leftGrid, nil
}

func generateWidgets(system model.Node) ([]container.Option, error) {
	builder := grid.New()
	el, elErr := createSystemElement(system)
	if elErr != nil {
		return nil, elErr
	}

	builder.Add(
		grid.ColWidthPerc(
			99,
			grid.Widget(
				el,
				container.BorderTitle(system.Name),
				container.Border(linestyle.Light),
			),
		),
	)

	gridOpts, builderErr := builder.Build()
	if builderErr != nil {
		return nil, fmt.Errorf("dashboard.Build => %v", builderErr)
	}

	return gridOpts, nil
}

func fakeStory(current float64) []float64 {
	var res []float64

	for i := 0; i < 100; i++ {
		v := math.Sin(float64(i) / 100 * math.Pi)
		res = append(res, v)
	}

	res = append(res, current)

	return res
}

func createSystemElement(sysNode model.Node) (widgetapi.Widget, error) {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
	)
	if err != nil {
		return nil, err
	}

	successLineErr := lc.Series(
		sysNode.Name,
		fakeStory(sysNode.RequestCount),
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.SeriesXLabels(map[int]string{
			0: "No requests",
		}),
	)
	if successLineErr != nil {
		return nil, successLineErr
	}

	errLineErr := lc.Series(
		sysNode.Name,
		fakeStory(sysNode.ErrorCount),
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.SeriesXLabels(map[int]string{
			0: "No errors",
		}),
	)

	if errLineErr != nil {
		return nil, successLineErr
	}

	return lc, nil
}