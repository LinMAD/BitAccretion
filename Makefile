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

## Compile fake provider
plugin_fake: golint gotest
	go build -buildmode=plugin -o ./provider.so plugin/fake/provider.go

## Compile new relic provider
plugin_relic: golint gotest
	go build -buildmode=plugin -o ./provider.so plugin/newrelic/provider.go
