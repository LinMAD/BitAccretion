#
# Makefile for BitAccretion

## Lint go files
golint:
	find . -type d -not -path "./.git/*" | xargs -L 1 golint

## Run tests for Go
gotest:
	- go test ./...
	- go test -race ./...

## Install project dependencies
prepare:
	go get -u golang.org/x/lint/golint
	go get -u github.com/golang/dep/cmd/dep
	go mod download

## Compile go code
build: golint gotest
	go build -o BitAccretion main.go

## Compile new relic processor
com_plugin_relic:
	go build -buildmode=plugin -o ./provider.so plugin/newrelic/provider.go
