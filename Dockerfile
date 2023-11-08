FROM golang:1.18.10-bullseye as builder

ENV DEV=/opt/kafka-gateway
RUN mkdir -p $DEV

WORKDIR $DEV
ADD . $DEV

RUN set -ex; \
  ls -ls; \
  go version; \
  go get -u github.com/confluentinc/confluent-kafka-go/v2/kafka; \
  go build -o ./bin/gateway;

FROM debian:latest
ENV DEV=/opt/kafka-gateway
RUN mkdir -p $DEV

WORKDIR $DEV
COPY entrypoint.sh entrypoint.sh
COPY --from=builder /opt/kafka-gateway/bin/gateway /opt/kafka-gateway/bin/gateway

ARG VCS=unknown
ARG BUILD_DATE=unknown

ENV PATH=$PATH:$DEV/bin

ENV BUILD_GIT_REF=${VCS}
ENV BUILD_DATE=${BUILD_DATE}

LABEL build.git.ref=${VCS} build.date=${BUILD_DATE} \
  org.opencontainers.image.source="https://github.com/sygnowski/kafka-gateway"

ENTRYPOINT ["sh", "/opt/kafka-gateway/entrypoint.sh"]
