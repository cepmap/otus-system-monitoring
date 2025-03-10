FROM golang:1.24 AS build

ENV BIN_FILE_SERVICE=/opt/stats-daemon/daemon
ENV BIN_FILE_CLIENT=/opt/stats-daemon/client
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE_SERVICE} cmd/daemon/* && \
    CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE_CLIENT} cmd/client/*

FROM ubuntu:latest

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="stats-daemon"
LABEL MAINTAINERS="cepmapp@gmail.com"

ENV BIN_FILE_SERVICE="/opt/stats-daemon/daemon"
ENV BIN_FILE_CLIENT="/opt/stats-daemon/client"

RUN apt-get update && apt-get install -y iftop sysstat

COPY --from=build ${BIN_FILE_SERVICE} ${BIN_FILE_SERVICE}
COPY --from=build ${BIN_FILE_CLIENT} ${BIN_FILE_CLIENT}

ENV CONFIG_FILE=/etc/stats-daemon/config.yaml
COPY ./_configs/config.yaml ${CONFIG_FILE}

CMD ["bash", "-c", "${BIN_FILE_SERVICE} --config ${CONFIG_FILE}"]