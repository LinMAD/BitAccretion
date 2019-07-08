#!/usr/bin/env sh

build() {
  docker build -t web_bit_accretion .
}

# ARG - 1 port
execute() {
  docker run -it --rm -p "80:8080" --device /dev/snd:/dev/snd --network host --name web_bit_accretion web_bit_accretion
}

### MAIN
build
execute "$1"
