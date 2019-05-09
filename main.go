package main

import (
	"encoding/json"
	"github.com/LinMAD/BitAccretion/core"
	"github.com/LinMAD/BitAccretion/core/api"
	"github.com/LinMAD/BitAccretion/core/cache"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"plugin"
)

const (
	// AppName human name
	AppName = "BitAccretion"
	// AppVersion human version id
	AppVersion = "0.7"
	tag        = "MAIN"
)

// appKernel container structure
type appKernel struct {
	// CacheManager key value storage
	CacheManager *cache.MemoryCache
	// router engine
	router *mux.Router
	// webRoot root path where stores static files: html, css
	webRoot string
	// monitoringProcessor controls traffic data preparation
	monitoringProcessor core.IProcessor
	// config is parsed app config file
	config Config
}

// Config contains configuration of application
type Config struct {
	GUI        bool   `json:"gui_monitoring"`
	Port       string `json:"web_port"`
	SurveyTime int    `json:"survey_time"`
}

// kernel stores main dependencies and configuration
var (
	kernel *appKernel
)

func init() {
	log.Printf("%s: Initializing %s - %s", tag, AppName, AppVersion)
	kernel = &appKernel{}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(tag + ": Could not retrieve working directory")
	}

	configPath := wd + "/config.json"
	kernel.webRoot = wd

	// Load kernel settings
	kernel.config = loadConfiguration(configPath)

	// Boot key value storage
	kernel.CacheManager = cache.Boot()

	// Create API and inject router engine
	kernel.router = mux.NewRouter()
	api.NewAPI(kernel.router, kernel.CacheManager, kernel.webRoot, kernel.config.SurveyTime).ServeAllRoutes(true)

	// Load processor plugin
	mod, err := plugin.Open(wd + "/processor.so")
	if err != nil {
		panic(err)
	}

	// Validate plugin - lookup for exported base structure
	prc, err := mod.Lookup("NewProcessor")
	if err != nil {
		log.Fatalf("Expected to be exported Processor structure in plugin, err: %v", err)
	}

	// Add implemented plugin to kernel
	kernel.monitoringProcessor = prc.(func() core.IProcessor)()
	// Load plugin settings
	kernel.monitoringProcessor.ParseConfig(configPath)
	kernel.monitoringProcessor.Prepare()
}

func main() {
	go func() {
		log.Printf("%s: Running web server at: http://127.0.0.1:%s/", tag, kernel.config.Port)
		log.Fatal(http.ListenAndServe(":"+kernel.config.Port, kernel.router))
	}()

	log.Printf("%s: Booting monitoring...", tag)

	for {
		graph := kernel.monitoringProcessor.GetLastAppGraph()

		kernel.CacheManager.Add(api.GraphStorageKey, graph)
		kernel.monitoringProcessor.Run()
	}
}

// loadConfiguration from json file
func loadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()

	if err != nil {
		log.Fatalln(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}
