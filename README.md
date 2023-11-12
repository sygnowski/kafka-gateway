# Kafka Http Gateway

Objective: Let's make HTTP flow supported via Kafka topic channel in order to support stream application.

### Configuration

Enviroment variables:
 - `BROKER`=kafka:9093
 - `PUBLISH_TOPIC`=gateway.publish
 - `SUBSCRIBE_TOPIC`=gateway.subscribe
 - `SUBSCRIBE_GROUP_ID`=gateway
 - `TIMEOUT`=30

### Commands

 - run `./dev.sh run host.docker.internal:9092`

### Tips

Setup the env `CGO_ENABLED=0` in order to build app with kafka lib.

Update notes:
```bash
# 2023-11-08

go get -u github.com/confluentinc/confluent-kafka-go/v2/kafka
go mod tidy

```

### Stack
 - [go-kafka](http://github.com/confluentinc/confluent-kafka-go)
 - [gin](https://github.com/gin-gonic/gin)
 - [uild](https://github.com/oklog/ulid)
 - [jwt](https://github.com/golang-jwt/jwt)
