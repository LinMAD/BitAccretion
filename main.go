// Copyright 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Binary donutdemo displays a couple of Donut widgets.
// Exist when 'q' is pressed.
package main

import (
	"context"
	"github.com/LinMAD/BitAccretion/dashboard"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/text"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/donut"
)

// playType indicates how to play a donut.
type playType int

const (
	playTypePercent playType = iota
	playTypeAbsolute
)

func main() {
	t, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())

	yellow, err := donut.New(
		donut.CellOpts(cell.FgColor(cell.ColorBlue)),
		donut.Label("System C", cell.FgColor(cell.ColorYellow)),
	)
	yellow.Percent(100)
	if err != nil {
		panic(err)
	}

	nodes := make([]model.Node, 4)
	nodes[0] = model.Node{
		Name:         "System 0",
		RequestCount: 101,
		ErrorCount:   232,
	}
	nodes[1] = model.Node{
		Name:         "System 1",
		RequestCount: 1102,
		ErrorCount:   104,
	}
	nodes[2] = model.Node{
		Name:         "System 2",
		RequestCount: 10,
		ErrorCount:   1,
	}
	nodes[3] = model.Node{
		Name:         "System 3",
		RequestCount: 2,
		ErrorCount:   0,
	}

	left, e := dashboard.CreateLeftLayout(nodes)
	if e != nil {
		panic(e)
	}

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.SplitVertical(
			left,
			// TODO Add right
			container.Right(
				container.Border(linestyle.Light),
				container.SplitHorizontal(
					container.Top(
						container.Border(linestyle.Light),
						EventLogWidget(),
					),
					container.Bottom(
						container.Border(linestyle.Light),
						container.PlaceWidget(yellow)),
				),
			),
		),
	)
	if err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(1*time.Second)); err != nil {
		panic(err)
	}
}

func EventLogWidget() container.Option {
	wrapped, err := text.New(text.WrapAtRunes())
	if err != nil {
		panic(err)
	}
	if err := wrapped.Write("2019:09:14 12:50| Error in system B \n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("2019:09:14 12:51| Error in system B \n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("2019:09:14 12:52| Error in system B \n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}
	if err := wrapped.Write("2019:09:14 12:53| Error in system C \n", text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		panic(err)
	}

	return container.PlaceWidget(wrapped)
}
