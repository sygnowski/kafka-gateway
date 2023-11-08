[![Total alerts](https://img.shields.io/lgtm/alerts/g/sygnowski/kafka-gateway.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/sygnowski/kafka-gateway/alerts/)

# Kafka Http Gateway

### Configuration

Enviroment variables:
 - `BROKER`=kafka:9093
 - `PUBLISH_TOPIC`=gateway.publish
 - `SUBSCRIBE_TOPIC`=gateway.subscribe
 - `SUBSCRIBE_GROUP_ID`=gateway
 - `TIMEOUT`=30

### Commands

 - run `./dev.sh run host.docker.internal:9092`


### Dependenices

[kafka support in go](http://github.com/confluentinc/confluent-kafka-go)

### Tips

Setup the env `CGO_ENABLED=0` in order to build app with kafka lib.

Update notes:
```bash
# 2023-11-08

go get -u github.com/confluentinc/confluent-kafka-go/v2/kafka
go mod tidy

```
