package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Properties struct {
	Server string
}

type KafkaProducer struct {
	prod *kafka.Producer
}

func (prod *KafkaProducer) Init(props *Properties) {
	fmt.Println("making kafka producer")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": props.Server})

	if err != nil {

		panic(err)
	}
	prod.prod = p

	go func() {
		for e := range p.Events() {
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
	}()

}
