#
# Makefile for BitAccretion

## Lint go files
golint:
	find . -type d -not -path "./.git/*" -not -path "./vendor/*" | xargs -L 1 golint

## Run tests for Go
gotest:
	- go test `go list ./... | grep -v vendor/ | grep -v extension/sound`
	- go test -race `go list ./... | grep -v vendor/ | grep -v extension/sound`

## Install project dependencies
prepare:
	go get -u golang.org/x/lint/golint
	go mod download

## Clean from artifacts
clean:
	rm -rf build && mkdir build
	mkdir -p build/resource

## Compile go code
build: golint gotest clean
	go build -o ./build/BitAccretion main.go

## Compile sound player
plugin_sound:
	rm -rf build/sound.so
	cp -r resource/sound ./build/resource
	go build -buildmode=plugin -o ./build/sound.so extension/sound/player.go

## Compile fake provider
plugin_fake:
	rm -rf build/provider.so
	go build -buildmode=plugin -o ./build/provider.so extension/fake/provider.go

## Compile new relic provider
plugin_relic:
	rm -rf build/provider.so
	go build -buildmode=plugin -o ./build/provider.so extension/newrelic/provider.go
