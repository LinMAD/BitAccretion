package main

import (
	"runtime"

	"github.com/LinMAD/BitAccretion/core"
	"github.com/LinMAD/BitAccretion/extension"
	"github.com/LinMAD/BitAccretion/model"
	"github.com/mum4k/termdash/terminal/termbox"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	t, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	kErr := warmKernel().Run(t)
	if kErr != nil {
		panic(kErr)
	}
}

// warmKernel for execution with loaded configuration and extensions
func warmKernel() *core.Kernel {
	var c *model.Config
	var p extension.IProvider
	var s extension.ISound
	var err error

	c, err = core.LoadConfig()
	if err != nil {
		panic(err)
	}

	p, err = core.LoadProviderPlugin()
	if err != nil {
		panic(err)
	}

	if c.IsSoundMode {
		s, err = core.LoadSoundPlugin()
		if err != nil {
			panic(err)
		}
	}

	return core.NewKernel(p, s, c)
}
