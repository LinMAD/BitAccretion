package main

import (
	"github.com/LinMAD/BitAccretion/core"
	"github.com/mum4k/termdash/terminal/termbox"
)

func main() {
	p, pErr := core.LoadProviderPlugin()
	if pErr != nil {
		panic(pErr)
	}

	t, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	kErr := core.NewKernel(p).Run(t)
	if kErr != nil {
		panic(kErr)
	}
}
