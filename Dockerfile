FROM golang:1.17.11-bullseye

RUN set -ex; \
  mkdir /opt/kafka-gateway; \
  go get -u gopkg.in/confluentinc/confluent-kafka-go.v1/kafka;

ARG VCS=unknown
ARG BUILD_DATE=unknown

ENV DEV=/opt/kafka-gateway
RUN mkdir -p $DEV

WORKDIR $DEV
ADD . $DEV

RUN set -ex; \
  ls -ls; \
  go version; \ 
  go build -o ./bin/gateway;

ENV PATH=$PATH:$DEV/bin

ENV BUILD_GIT_REF=${VCS}
ENV BUILD_DATE=${BUILD_DATE}

LABEL build.git.ref=${VCS} build.date=${BUILD_DATE}

ENTRYPOINT ["sh", "/opt/kafka-gateway/entrypoint.sh"]
