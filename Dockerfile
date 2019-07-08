FROM golang:1.12-stretch

LABEL description="Image to execute BitAccretion in web"
ENV GO111MODULE=off

# Install deps
RUN apt-get -y update && apt-get install -y \
    ca-certificates \
    make \
    libasound2 \
    libasound2-dev \
    curl \
    tar

# Install gotty for web terminal
# TODO Remove hardocded gotty version
RUN curl -sLk https://github.com/yudai/gotty/releases/download/v2.0.0-alpha.3/gotty_2.0.0-alpha.3_linux_amd64.tar.gz | tar xzC /opt && \
    chmod  +x /opt/gotty && \
    chmod 744 /opt/gotty && \
    mv /opt/gotty /usr/local/bin/

RUN apt-get purge --auto-remove -y && apt-get clean && rm -rf /var/lib/apt/lists*

# Setup application paths
RUN mkdir -p /go/src/github.com/LinMAD/BitAccretion
COPY ./. /go/src/github.com/LinMAD/BitAccretion
WORKDIR /go/src/github.com/LinMAD/BitAccretion

# Compile it
# TODO Remove hardcoded provider plugin
RUN go build -o ./build/BitAccretion main.go && make plugin_relic && make plugin_sound
COPY config.json build/config.json
WORKDIR build

# Execute
EXPOSE 8080
CMD gotty --reconnect ./BitAccretion
