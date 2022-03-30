echo Gateway Entrypoint
echo
echo BROKER: $BROKER
echo PUBLISH_TOPIC: $PUBLISH_TOPIC
echo SUBSCRIBE_TOPIC: $SUBSCRIBE_TOPIC
echo SUBSCRIBE_GROUP_ID: $SUBSCRIBE_GROUP_ID
echo TIMEOUT: $TIMEOUT
echo

ls -la /app

/app/gateway
