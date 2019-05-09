## BitAccretion

[![Build Status](https://travis-ci.org/LinMAD/BitAccretion.svg?branch=master)](https://travis-ci.org/LinMAD/BitAccretion)

[Project plot on medium](https://medium.com/@artjomnemiro/how-valuable-can-be-visual-monitoring-923e9e865625)

![](./resources/Sample.png)

### Project structure
```text
. BitAccretion
├── build           // Compiled and packed files, ready for deploying
├── vendor          // External packages, dependencies 
├── public          // Compiled static files from the react app
├── resources       // UI resources 
│   └── components  // React components
├── plugins         // Processor plugins, basic implementation to construct graph with data
└── core            // Core system files
```

### How to prepare to develop or compile
Easy - use docker or install all OS dependencies manually, but...
##### !!! Compile only on Linux, the issue with Go Plugins on others OS !!!

Docker:
- First, build docker image where you will build all code.
- Second, compile all needed parts like: `js, go plugins and go app`

Example with docker
```text
$: docker build . -t bitaccretion:latest
$: docker run --volume `pwd`:/go/src/github.com/LinMAD/BitAccretion --name bit_build --rm -it bitaccretion:latest /bin/sh
```

To install project dependencies run:
```text
$: make prepare
```
It's install node modules and external Golang packages

To prepare static data like: JS, CSS, HTML
```text
$: make js
```
It compiles React to static files.

Example how to compile the plugin
```text
$: make plugin_relic
```

Build project for deployment, all files will be in `build` folders for "production"
Example how executing the build 
```text
$: make build
```

Compiled folder can be look like that:
```text
├── BitAccretion
├── config.json
├── processor.so
├── resources
└── sound.so
```

### Implementing own plugin

To implement own processor for Netflix's Vizceral you can write a plugin.
Each plugin implements interface from package `core` file `processorInterface.go` interface `IProcessor`

One tricky moment with plugins.
I didn't found an elegant way to export function to create a pointer to the processor, so you must do it like that:
```go
// Return an interface but create a pointer to structure with implemented methods
func NewProcessor() core.IProcessor {
    return new(YourProcessorStructure)
}

```

Anyway, as for example, you can take already implemented plugins.

P.S. If you will compile sound plugin then install dependency for `libsamplerate`.
Example: ```sudo apt install libsamplerate0```

## Configuration
There is a config file `config.json`.

```text
{
  "survey_time": 1,                // In seconds
  "web_port": "8080",
  "api_key": "your_key_if_needed", // Depending on plugin what must be surveyed
  "enabled_sound_alert": false,
  "health_sensitivity": {          // Settings of conversion rate, to mark warning, alert node
    "danger": 10,
    "warning": 0
  },
  "app_sets": [
    {
      "name": "appName",
      "id": "newRelicId",
      "nested": [
        {
          "name": "relatedApp",
          "id": "newRelicId",
          "relic_metrics": [
            "HttpDispatcher",
            "Errors/all"
          ]
        }
      ],
      "relic_metrics": [
        "HttpDispatcher",
        "Errors/all"
      ]
    }
  ]
}
```

TODO List
--------------
 - Add flag to disable\enable logs
 - Return error to main loop (No fatal, panic in plugins and core)
