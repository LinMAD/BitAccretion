package dashboard

import (
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/text"
)

func NewDashboardContainer(t terminalapi.Terminal) (*container.Container, error) {
	left, e := CreateLeftLayout(GetStubNodes())
	if e != nil {
		panic(e)
	}

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.SplitVertical(
			left,
			container.Right(
				container.SplitHorizontal(
					container.Top(
						container.Border(linestyle.Light),
						EventLogWidget(),
					),
					container.Bottom(
						container.Border(linestyle.Light),
					),
				),
			),
		),
	)

	return c, err
}

func EventLogWidget() container.Option {
	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:50| Error in CSEye\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:51| Error in CSEye\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:52| Error in CSEye\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("|2019:09:14 12:53| Error in Fastlane\n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}

	return container.PlaceWidget(wrapped)
}
