FROM golang:1.17.11-bullseye as stack

RUN set -ex; \
  mkdir /opt/kafka-gateway; \
  go get -u gopkg.in/confluentinc/confluent-kafka-go.v1/kafka;

FROM stack

COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["sh", "/entrypoint.sh"]

ENV DEV=/opt/kafka-gateway
RUN mkdir -p $DEV

WORKDIR $DEV
ADD . $DEV

RUN set -ex; \
  ls -ls; \
  go version; \ 
  go build -o ./bin/gateway;
