package kafint

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const CORRELATION = "Correlation"

type Properties struct {
	Server           string
	PublishTopic     string
	SubscribeTopic   string
	SubscribeGroupId string
	Timeout          uint32
}

type Correlaction struct {
	id   string
	resp chan kafka.Message
}

type KafkaIntegrator struct {
	prod           *kafka.Producer
	cons           *kafka.Consumer
	cfg            *Properties
	correlationMap sync.Map
	timeout        time.Duration
}

func (this *KafkaIntegrator) Init(props *Properties) {
	this.cfg = props
	this.timeout = time.Duration(props.Timeout) * time.Second
	fmt.Printf("Timeout :%s\n", this.timeout)

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
	if err != nil {
		panic(err)
	}
	this.cons = c

	this.correlationMap = sync.Map{}

	c.Subscribe(props.SubscribeTopic, nil)
	go this.fetch()
}

func (kint *KafkaIntegrator) fetch() {
	run := true

	for run {
		select {
		// case sig := <-sigchan:
		// 	fmt.Printf("Caught signal %v: terminating\n", sig)
		// 	run = false
		default:
			ev := kint.cons.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
				kint.onNewMessage(e)
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}

func (kint *KafkaIntegrator) onNewMessage(m *kafka.Message) {
	if m.Headers != nil {
		for _, h := range m.Headers {
			if h.Key == CORRELATION {
				cid := string(h.Value)
				println("got cid: " + cid)

				c, _ := kint.correlationMap.Load(cid)
				if c != nil {
					corr := c.(*Correlaction)
					corr.resp <- *m
					println("cid sent")
					break
				} else {
					println("no map hit for " + cid)
				}
			}
		}
	}
}

func (kint *KafkaIntegrator) Publish(w http.ResponseWriter, req *http.Request) {
	var bodyBytes []byte
	var err error

	cid := req.Header.Get(CORRELATION)
	if cid == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Missing the Correlation Header.")
	} else {

		if req.Body != nil {
			bodyBytes, err = io.ReadAll(req.Body)
			if err != nil {
				fmt.Printf("Body reading error: %v", err)
				return
			}
			defer req.Body.Close()
		}

		kint.publishToKafka(bodyBytes, cid)

		c := &Correlaction{id: cid, resp: make(chan kafka.Message, 1)}
		kint.correlationMap.Store(cid, c)

		clean := func() {
			close(c.resp)
			kint.correlationMap.Delete(cid)
		}

		select {
		case data := <-c.resp:
			w.Write(data.Value)
			clean()
			break
		case <-time.After(kint.timeout):
			w.WriteHeader(http.StatusRequestTimeout)
			io.WriteString(w, "Timeout")
			clean()
		}
	}
}

func (prod *KafkaIntegrator) publishToKafka(data []byte, cid string) {
	fmt.Println("Publishing to Kafka...")
	fmt.Println(string(data))

	deliveryChan := make(chan kafka.Event)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &prod.cfg.PublishTopic,
			Partition: kafka.PartitionAny},
		Value:   data,
		Headers: []kafka.Header{kafka.Header{Key: CORRELATION, Value: []byte(cid)}},
	}

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
