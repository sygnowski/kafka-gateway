#!/bin/bash

main () {

    case $1 in
        run)
            dockerRun $2
            ;;
        build)
            dockerBuild
            ;;
        *)
        badOpt $@
    esac

}

badOpt() {
    echo bad options: $@
}

dockerRun() {
    docker run \
     --name kafka-gateway \
     -p 8080:8080 \
     -e BROKER=$1 \
     -e PUBLISH_TOPIC=gateway.publish \
     -e SUBSCRIBE_TOPIC=gateway.subscribe \
     -e SUBSCRIBE_GROUP_ID=gateway \
     --network dev-network \
     -it --rm s7i/kafka-gateway
}

dockerBuild() {
    docker build -t s7i/kafka-gateway .
}

main $@