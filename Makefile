#
# Makefile for BitAccretion

## Lint go files
golint:
	find . -type d -not -path "./.git/*" -not -path "./vendor/*" -not -path "./node_modules/*" | xargs -L 1 golint

## Run tests for Go
gotest:
	go test ./core/...
	go test ./

## Clean compiled react build
clean:
	rm -rf build && mkdir build
	rm -rf public && mkdir public
	mkdir -p build/resources

## Install project dependencies
prepare:
	go get -u golang.org/x/lint/golint
	go get -u github.com/golang/dep/cmd/dep
	npm install
	dep ensure

## Prepare for development work-flow react app
js: clean
	npm run build
	cp -rl  public/* build/resources
	rm -rf public

## Compile core go app
com_core:
	go build -o BitAccretion main.go && mv BitAccretion build/.

## Compile new relic processor
com_plugin_relic:
	rm -rf build/processor.so
	go build -buildmode=plugin -o ./build/processor.so plugins/relic/newrelic.go

## Compile new relic processor
com_plugin_sound:
	rm -rf build/resources/sound && rm -rf build/sound.so && mkdir -p build/resources/sound
	cp -r resources/sound build/resources
	go build -buildmode=plugin -o ./build/sound.so plugins/misc/sound.go

## Build go code
build_go: com_core com_plugin_relic com_plugin_sound

## Build all parts
build_full: js build_go com_plugin_relic com_plugin_sound