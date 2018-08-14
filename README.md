## BitAccretion

##### !!! Compile only on Linux, the issue with Go Plugins !!!

[![Build Status](https://travis-ci.org/LinMAD/BitAccretion.svg?branch=master)](https://travis-ci.org/LinMAD/BitAccretion)

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
To install project dependencies run:
```text
$: make prepare
```
It's install node modules and external Golang packages

To prepare static data like: JS, CSS, HTML
```text
$: make js
```
It compiles React to static files and creates bind data to serve it via HTTP.

Example how to compile the plugin
```text
$: make plugin_relic
```

Build project for deployment, all files will be in `build` folders for "production"
Example how executing the build 
```text
$: make build
```

Example of compiled folder:
```text
├── BitAccretion
├── config.json
├── processor.so
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

```json
{
  "web_port": "8080",
  "api_key": "your_key_if_needed",
  "health_sensitivity": {
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
