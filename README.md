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
