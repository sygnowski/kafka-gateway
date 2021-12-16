FROM golang:latest as stack

RUN set -ex; \
  mkdir /opt/kafka-gateway; \
  go get -u gopkg.in/confluentinc/confluent-kafka-go.v1/kafka;

FROM stack

WORKDIR /opt/kafka-gateway
COPY go.mod go.mod
COPY go.sum go.sum
COPY internal internal
COPY cmd cmd
COPY main.go main.go

RUN set -ex; \
  ls -ls; \
  go version; \ 
  go build -o ./bin/gateway;