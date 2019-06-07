package main

import (
	"log"
	"os"
	"plugin"

	"github.com/LinMAD/BitAccretion/kernel"
	"github.com/LinMAD/BitAccretion/provider"
	"github.com/mum4k/termdash/terminal/termbox"
)

// TODO Refactor main, configuration and clean up
var (
	providerImpl provider.IProvider
	configPath   string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic("Could not retrieve working directory, error: " + err.Error())
	}
	configPath = wd + "/config.json"

	mod, err := plugin.Open(wd + "/provider.so")
	if err != nil {
		panic("Unable to open provider.so plugin, error: " + err.Error())
	}

	// Validate plugin - lookup for exported base function to get implementation
	prc, err := mod.Lookup("NewProvider")
	if err != nil {
		log.Fatalf("Expected to be exported Processor structure in plugin, err: %v", err)
	}

	// Add implemented plugin to kernel
	providerImpl = prc.(func() provider.IProvider)()
}

func main() {
	t, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	k := kernel.NewKernel(providerImpl)
	kErr := k.Run(t)
	if kErr != nil {
		panic(kErr)
	}
}
