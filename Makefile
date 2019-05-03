#
# Makefile for BitAccretion

## Lint go files
golint:
	find . -type d -not -path "./.git/*" -not -path "./vendor/*" -not -path "./node_modules/*" | xargs -L 1 golint

## Run tests for Go
gotest:
	go test ./core/...
	go test ./

## Prepare for development work-flow react app
js: clean
	npm run build
	mkdir -p build/resources
	cp -rl  public/* build/resources

## Clean compiled react build
clean:
	rm -rf build && mkdir build
	rm -rf public && mkdir public

## Build go code
build_go: plugin_relic plugin_sound
	go build -o BitAccretion main.go && mv BitAccretion build/.

## Build all parts
build_full: js build_go plugin_relic plugin_sound

## Compile new relic processor
plugin_relic:
	rm -rf build/processor.so
	go build -buildmode=plugin -o ./build/processor.so plugins/relic/newrelic.go

## Compile new relic processor
plugin_sound:
	rm -rf build/resources/sound && rm -rf build/sound.so && mkdir -p build/resources/sound
	cp -r resources/sound build/resources
	go build -buildmode=plugin -o ./build/sound.so plugins/misc/sound.go

## Install project dependencies
prepare:
	go get -u golang.org/x/lint/golint
	go get -u github.com/golang/dep/cmd/dep
	npm install
	dep ensure
