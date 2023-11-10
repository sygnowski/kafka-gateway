#!/bin/bash

#set -ex

function info() {
    cat <<EOF
Gateway Entrypoint

BROKER: $BROKER
PUBLISH_TOPIC: $PUBLISH_TOPIC
SUBSCRIBE_TOPIC: $SUBSCRIBE_TOPIC
SUBSCRIBE_GROUP_ID: $SUBSCRIBE_GROUP_ID
TIMEOUT: $TIMEOUT

BUILD_GIT_REF: $BUILD_GIT_REF
BUILD_DATE: $BUILD_DATE

YAML_CONFIG: $APP_CONFIG
EOF

}

function asYaml() {
    cat <<EOF
app:
  port: 8080
  timeout: ${TIMEOUT}
  jwt-chk: false
  id-gen-kind: ulid
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

export APP_CONFIG=${APP_CONFIG:-$(asYaml)}
info

exec gateway
