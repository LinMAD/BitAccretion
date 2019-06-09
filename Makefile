#
# Makefile for BitAccretion

## Lint go files
golint:
	find . -type d -not -path "./.git/*" -not -path "./vendor/*" | xargs -L 1 golint

## Run tests for Go
gotest:
	- go test `go list ./... | grep -v /vendor/`
	- go test -race `go list ./... | grep -v /vendor/`

## Install project dependencies
prepare:
	go get -u golang.org/x/lint/golint
	go mod download

## Compile go code
build: golint gotest
	go build -o BitAccretion main.go

## Compile sound player
plugin_sound:
	go build -buildmode=plugin -o ./sound.so extension/sound/player.go

## Compile fake provider
plugin_fake:
	go build -buildmode=plugin -o ./provider.so extension/fake/provider.go

## Compile new relic provider
plugin_relic:
	go build -buildmode=plugin -o ./provider.so extension/newrelic/provider.go
