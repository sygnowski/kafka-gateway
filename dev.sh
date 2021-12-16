#!/bin/bash

main () {

    case $1 in
        run)
            dockerRun
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
    docker run --name kafka-gateway -it --rm s7i/kafka-gateway    
}

dockerBuild() {
    docker build -t s7i/kafka-gateway .
}

main $@