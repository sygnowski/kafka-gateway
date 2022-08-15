export BROKER="localhost:9092"
export PUBLISH_TOPIC="TodoAction"
export SUBSCRIBE_TOPIC="TodoReaction"
export SUBSCRIBE_GROUP_ID="todo-gateway"
export TIMEOUT=60

./kafka-gateway

