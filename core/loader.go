package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"plugin"

	"github.com/LinMAD/BitAccretion/model"
	"github.com/LinMAD/BitAccretion/provider"
)

const (
	configFile = "config.json"

	pluginProviderFile         = "provider.so"
	pluginProviderMainFunction = "NewProvider"

	pluginSoundFile         = "sound.so"
	pluginSoundMainFunction = "PlayAlert"
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

// LoadConfig resolve core settings
func LoadConfig() (*model.Config, error) {
	var config model.Config
	c, cErr := os.Open(path.Join(wd, configFile))
	if cErr != nil {
		return nil, fmt.Errorf(cErr.Error())
	}
	defer c.Close()

	jsonParser := json.NewDecoder(c)
	parseErr := jsonParser.Decode(&config)
	if parseErr != nil {
		return nil, fmt.Errorf(parseErr.Error())
	}

	return &config, nil
}

// LoadProviderPlugin resolve data provider plugin
func LoadProviderPlugin() (provider.IProvider, error) {
	mod, modErr := plugin.Open(path.Join(wd, pluginProviderFile))
	if modErr != nil {
		return nil, fmt.Errorf("unable to open plugin %s, error: %s", pluginProviderFile, modErr.Error())
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

// LoadSoundPlugin resolve sound plugin
func LoadSoundPlugin() (func(name string), error) {
	mod, modErr := plugin.Open(path.Join(wd, pluginSoundFile))
	if modErr != nil {
		return nil, fmt.Errorf("unable to open plugin %s, error: %s", pluginSoundFile, modErr.Error())
	}

	s, sErr := mod.Lookup(pluginSoundMainFunction)
	if sErr != nil {
		return nil, fmt.Errorf("expected to be found '%s' function in plugin, err: %s", pluginSoundMainFunction, sErr.Error())
	}

	return s.(func(name string)), nil
}
