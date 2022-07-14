FROM golang:1.17.11-bullseye as builder

ENV DEV=/opt/kafka-gateway

WORKDIR $DEV
ADD . $DEV

RUN set -ex; \
  mkdir -p $DEV; \
  ls -ls; \
  go version; \ 
  go get -tags musl -u github.com/confluentinc/confluent-kafka-go/kafka; \
  go build -tags musl -o ./bin/gateway;

FROM alpine:latest

ARG VCS=unknown
ARG BUILD_DATE=unknown

ENV PUBLISH_TOPIC=gateway.publish \
  SUBSCRIBE_TOPIC=gateway.subscribe \
  SUBSCRIBE_GROUP_ID=gateway \
  TIMEOUT=30

ENV DEV=/opt/kafka-gateway
RUN mkdir -p $DEV
ENV PATH=$PATH:$DEV/bin

WORKDIR $DEV

ENV BUILD_GIT_REF=${VCS}
ENV BUILD_DATE=${BUILD_DATE}

COPY --from=builder /opt/kafka-gateway/bin/gateway /opt/kafka-gateway/bin/gateway
COPY entrypoint.sh /opt/kafka-gateway/entrypoint.sh

LABEL build.git.ref=${VCS} build.date=${BUILD_DATE}

ENTRYPOINT ["sh", "/opt/kafka-gateway/entrypoint.sh"]
