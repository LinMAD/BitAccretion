package main

import (
	"context"
	"time"

	"github.com/LinMAD/BitAccretion/dashboard"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

func main() {
	t, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	// TODO Move that to dashboard file -> setup, observer updates

	ctx, cancel := context.WithCancel(context.Background())
	c, err := dashboard.NewDashboardContainer(t)
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

