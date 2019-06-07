package core

import (
	"fmt"
	"github.com/LinMAD/BitAccretion/provider"
	"os"
	"path"
	"plugin"
)

const (
	config                     = "config.json"
	pluginProvider             = "provider.so"
	pluginProviderMainFunction = "NewProvider"
)

// wd resolved rooted path name corresponding to the current directory
var wd string

func init() {
	if wd != "" {
		return
	}

	var err error
	wd, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("could not retrieve working directory, error: %s", err.Error()))
	}
}

// LoadProviderPlugin resolve data provider plugin
func LoadProviderPlugin() (provider.IProvider, error) {
	mod, modErr := plugin.Open(path.Join(wd, pluginProvider))
	if modErr != nil {
		return nil, fmt.Errorf("unable to open plugin provider.so, error: " + modErr.Error())
	}

	p, pErr := mod.Lookup(pluginProviderMainFunction)
	if pErr != nil {
		return nil, fmt.Errorf(
			"expected to be found '%s' function in plugin, err: %s",
			pluginProviderMainFunction,
			pErr.Error(),
		)
	}

	return p.(func() provider.IProvider)(), nil
}
