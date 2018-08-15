FROM golang:1.10-alpine

LABEL description="BitAccretion - dockerized container for compiling"

# Create working paths
RUN mkdir -p /go/src/github.com/LinMAD/BitAccretion
WORKDIR /go/src/github.com/LinMAD/BitAccretion

# Install dependecies
RUN apk update && apk add --update \
    build-base make \
    git \
    nodejs nodejs-npm \
    alsa-lib-dev

