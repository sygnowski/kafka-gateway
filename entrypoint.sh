#!/bin/bash

#set -ex

function info() {
    echo Gateway Entrypoint
    echo
    echo BROKER: $BROKER
    echo PUBLISH_TOPIC: $PUBLISH_TOPIC
    echo SUBSCRIBE_TOPIC: $SUBSCRIBE_TOPIC
    echo SUBSCRIBE_GROUP_ID: $SUBSCRIBE_GROUP_ID
    echo TIMEOUT: $TIMEOUT
    echo
    echo BUILD_GIT_REF: $BUILD_GIT_REF
    echo BUILD_DATE: $BUILD_DATE
    echo
}

function asYaml() {
    cat <<EOF
app:
  port: 8080
  timeout: ${TIMEOUT}
  jwt-chk: false
  id-gen-type: uuid
  context: /
pub:
  topic: ${PUBLISH_TOPIC}
  properties:
   - "bootstrap.servers=${BROKER}"
sub:
  topic: ${SUBSCRIBE_TOPIC}
  properties:
   - "bootstrap.servers=${BROKER}"
   - "group.id=${SUBSCRIBE_GROUP_ID}"
   - "auto.offset.reset=earliest"

EOF
}

info

export APP_CONFIG=$(asYaml)

exec gateway
