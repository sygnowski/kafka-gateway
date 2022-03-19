package kafint

import (
	"fmt"
	"io"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Properties struct {
	Server           string
	PublishTopic     string
	SubscribeTopic   string
	SubscribeGroupId string
}

type KafkaIntegrator struct {
	prod *kafka.Producer
	cfg  *Properties
}

func (this *KafkaIntegrator) Init(props *Properties) {
	this.cfg = props

	fmt.Println("making kafka producer")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": props.Server})

	if err != nil {

		panic(err)
	}
	this.prod = p

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": props.Server,
		"group.id":          props.SubscribeGroupId,
		"auto.offset.reset": "latest",
	})

	c.Subscribe(props.SubscribeTopic, nil)

}

func (prod *KafkaIntegrator) Publish(w http.ResponseWriter, req *http.Request) {
	var bodyBytes []byte
	var err error

	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer req.Body.Close()
	}
	prod.publishToKafka(bodyBytes)

}

func (prod *KafkaIntegrator) publishToKafka(data []byte) {
	fmt.Println("Publishing to Kafka...")
	fmt.Println(string(data))

	deliveryChan := make(chan kafka.Event)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &prod.cfg.PublishTopic,
			Partition: kafka.PartitionAny},
		Value: data}

	prod.prod.Produce(msg, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)
	onMessage(m)

	close(deliveryChan)
}

func handleEvents(events chan kafka.Event) {

	for e := range events {
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			onMessage(m)
			return

		default:
			fmt.Printf("Ignored event: %s\n", ev)
		}
	}
}

func onMessage(m *kafka.Message) {
	fmt.Println(m)
	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
}
