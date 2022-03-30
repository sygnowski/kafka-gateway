FROM golang:latest as stack

RUN set -ex; \
  mkdir /opt/kafka-gateway; \
  go get -u gopkg.in/confluentinc/confluent-kafka-go.v1/kafka;

FROM stack as builder

ENV DEV=/opt/kafka-gateway CGO_ENABLED=0

WORKDIR $DEV
ADD . $DEV

RUN set -ex; \
  ls -ls; \
  go version; \ 
  go build -tags netgo -a -v -o ./bin/gateway;

FROM alpine:latest
ENV PUBLISH_TOPIC=gateway.publish \
  SUBSCRIBE_TOPIC=gateway.subscribe \
  SUBSCRIBE_GROUP_ID=gateway \
  TIMEOUT=30

COPY --from=builder /opt/kafka-gateway/bin/gateway /app/gateway
COPY entrypoint.sh /app/entrypoint.sh
ENTRYPOINT /app/entrypoint.sh
EXPOSE 8080