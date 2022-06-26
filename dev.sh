#!/bin/bash
 
 PUBLISH_TOPIC=gateway.publish
 SUBSCRIBE_TOPIC=gateway.subscribe
 SUBSCRIBE_GROUP_ID=gateway
 TIMEOUT=30
 
 DOCKER_IMAGE=s7i/kafka-gateway
 DOCKER_NETWORK=dev-network

main () {

    case $1 in
        run)
            run $2
            ;;
        dr)
            dockerRun $2
            ;;
        db)
            dockerBuild
            ;;
        win)
            buildWindwsBins
            ;;
        test)
            runTest
            ;;
        *)
        badOpt $@
    esac

}

badOpt() {
    echo bad options: $@
}

run() {
    export BROKER=$1
    export PUBLISH_TOPIC=$PUBLISH_TOPIC
    export SUBSCRIBE_TOPIC=$SUBSCRIBE_TOPIC
    export SUBSCRIBE_GROUP_ID=$SUBSCRIBE_GROUP_ID
    export TIMEOUT=$TIMEOUT
    ./kafka-gateway
}

dockerRun() {
    docker run \
     --name kafka-gateway \
     -p 8080:8080 \
     -e BROKER=$1 \
     -e PUBLISH_TOPIC=$PUBLISH_TOPIC \
     -e SUBSCRIBE_TOPIC=$SUBSCRIBE_TOPIC \
     -e SUBSCRIBE_GROUP_ID=$SUBSCRIBE_GROUP_ID \
     -e TIMEOUT=$TIMEOUT \
     --network $DOCKER_NETWORK\
     -it --rm $DOCKER_IMAGE
}

dockerBuild() {
    docker build -t $DOCKER_IMAGE .
}

buildWindwsBins() {
    GOOS=windows GOARCH=amd64 go build -o bin/gateway-amd64.exe
}

runTest() {
    go test -v -timeout 30s s7i.io/kafka-gateway/internal/kafint s7i.io/kafka-gateway/internal/util
}

main $@