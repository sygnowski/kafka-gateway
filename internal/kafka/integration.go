package kafka

import (
	"fmt"
	"io"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Properties struct {
	Server       string
	PublishTopic string
}

type KafkaProducer struct {
	prod            *kafka.Producer
	cfg             *Properties
	producerChannel chan *kafka.Message
}

func (this *KafkaProducer) Init(props *Properties) {
	fmt.Println("making kafka producer")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": props.Server})

	if err != nil {

		panic(err)
	}
	this.prod = p
	this.cfg = props

	this.producerChannel = p.ProduceChannel()

	go this.handleEvents(p.Events())

}

func (prod *KafkaProducer) Publish(w http.ResponseWriter, req *http.Request) {
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

func (prod *KafkaProducer) publishToKafka(data []byte) {
	fmt.Println("Publishing to Kafka...")
	fmt.Println(string(data))

	prod.producerChannel <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &prod.cfg.PublishTopic,
			Partition: 0},
		Value: data}
}

func (prod *KafkaProducer) handleEvents(events chan kafka.Event) {

	for e := range events {
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			if m.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
			return

		default:
			fmt.Printf("Ignored event: %s\n", ev)
		}
	}
}
