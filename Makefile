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
	go-bindata-assetfs -pkg api -o core/api/static_bin_data.go public/...

## Clean compiled react build
clean:
	npm run clean

## Build application
build: js
	rm -rf build && mkdir build
	go build -o BitAccretion main.go && mv BitAccretion build/.
	cp config.json build/.

## Compile new relic processor
plugin_relic:
	go build -buildmode=plugin -o ./build/processor.so plugins/relic/newrelic.go

## Install project dependencies
prepare:
	go get -u golang.org/x/lint/golint
	go get -u github.com/go-bindata/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...
	npm install
	dep ensure
