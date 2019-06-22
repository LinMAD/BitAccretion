## BitAccretion

[![Build Status](https://travis-ci.org/LinMAD/BitAccretion.svg?branch=master)](https://travis-ci.org/LinMAD/BitAccretion)

#### About
`BitAccretion` - simple tool to visualize system metrics.
Metrics are collected by the provider and it's based on Go plugins to support different ways to aggregate\assemble metrics data.

[More why it was created, post on medium](https://medium.com/@artjomnemiro/how-valuable-can-be-visual-monitoring-923e9e865625)

##### Structure

[Configuration file](./config.json.tpl)

Folder hierarchy:
```text
BitAccretion
├── build           // Stores compiled project (if used make)
├── core            // Responsible for communication between dashboard and provider
├── dashboard       // Terminal UI with logic handling
├── event           // Observer with common interfaces
├── extension       // Go plugin folder
│   ├── fake
│   ├── newrelic
│   ├── sound
|   ├── provider.go // Provider (metrics assembling) interface for plugin
|   └── sound.go    // Sound (alerting) interface for plugin
├── logger          // Application logs interface
├── model           // Data structures
├── resource        // External files, images, audio files, etc
│   └── sound
│       ├── alarm
│       └── voice
├── util            // Small useful functions
└── vendor          // Third party dependencies (managed by go modules)
```

#### Example how it's looks like
![Demo example](./resource/example.gif)

##### TODO List:
```text
TODO GOTTY in docker container to access from browser
TODO Remove clock and add degradation chart from exec time (Show error regression from exec time till now)
TODO Add to dashboard name processor name (to displace source of data)
```