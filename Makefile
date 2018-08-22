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
	mv -f public/* build/resources

## Clean compiled react build
clean:
	rm -rf public && mkdir public

## Build application
build: js
	rm -rf build && mkdir build
	go build -o BitAccretion main.go && mv BitAccretion build/.

## Compile new relic processor
plugin_relic:
	go build -buildmode=plugin -o ./build/processor.so plugins/relic/newrelic.go

## Compile new relic processor
plugin_sound:
	go build -buildmode=plugin -o ./build/sound.so plugins/misc/sound.go
	rm -rf build/resources/sound && mkdir -p build/resources/sound
	cp -r resources/sound build/resources

## Install project dependencies
prepare:
	go get -u golang.org/x/lint/golint
	go get -u github.com/golang/dep/cmd/dep
	npm install
	dep ensure
