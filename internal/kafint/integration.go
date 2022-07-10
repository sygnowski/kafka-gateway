package kafint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"s7i.io/kafka-gateway/internal/util"
)

const CORRELATION string = "correlation"
const CONTEXT string = "context"

var correlationPath = []string{CONTEXT, CORRELATION}

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

func NewKafkaIntegrator(props *Properties) *KafkaIntegrator {
	kint := KafkaIntegrator{}
	kint.init(props)
	return &kint
}

func (ki *KafkaIntegrator) init(props *Properties) {
	ki.cfg = props
	ki.timeout = time.Duration(props.Timeout) * time.Second
	fmt.Printf("Timeout :%s\n", ki.timeout)

	fmt.Println("making kafka producer")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": props.Server})

	if err != nil {

		panic(err)
	}
	ki.prod = p

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": props.Server,
		"group.id":          props.SubscribeGroupId,
		"auto.offset.reset": "latest",
	})
	if err != nil {
		panic(err)
	}
	ki.cons = c

	ki.correlationMap = sync.Map{}

	c.Subscribe(props.SubscribeTopic, nil)
	go ki.fetch()
}

func (ki *KafkaIntegrator) fetch() {
	run := true

	for run {
		select {
		default:
			ev := ki.cons.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("[KAFKA] Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("[KAFKA] Headers: %v\n", e.Headers)
				}
				ki.onNewMessage(e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "[KAFKA ERROR] %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					waitTime := 10 * time.Second
					fmt.Fprintf(os.Stderr, "[KAFKA] Waiting for broker... %s\n", waitTime)
					time.Sleep(waitTime)
				}
			default:
				fmt.Printf("[KAFKA IGNORED EVENT] %v\n", e)
			}
		}
	}
}

func (ki *KafkaIntegrator) onNewMessage(m *kafka.Message) {
	resp := make(map[string]interface{})

	if err := json.Unmarshal(m.Value, &resp); err == nil {
		if c := resp[CONTEXT]; c != nil {
			ctx := c.(map[string]interface{})

			if cid := ctx[CORRELATION]; cid != nil {
				if c, _ := ki.correlationMap.Load(cid); c != nil {
					println("got cid in the body.context.correlation: " + cid.(string))
					corr := c.(*Correlaction)
					corr.resp <- *m
					return
				}
			}
		}
	}

	if m.Headers != nil {
		for _, h := range m.Headers {
			if h.Key == CORRELATION {
				cid := string(h.Value)
				println("got cid: " + cid)

				c, _ := ki.correlationMap.Load(cid)
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

func (ki *KafkaIntegrator) Publish(w http.ResponseWriter, req *http.Request) {
	var bodyBytes []byte
	var err error
	var ctxAtt bool

	cid := req.Header.Get(CORRELATION)
	if cid == "" {
		cid = util.UUID()
	}

	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer req.Body.Close()
	}
	ctxAtt, bodyBytes = ki.attachContext(bodyBytes, cid)
	if ctxAtt {
		fmt.Printf("[KI] New context attached [%s].\n", cid)
	}
	ki.publishToKafka(bodyBytes, cid)

	c := &Correlaction{id: cid, resp: make(chan kafka.Message, 1)}
	ki.correlationMap.Store(cid, c)

	clean := func() {
		close(c.resp)
		ki.correlationMap.Delete(cid)
	}

	select {
	case data := <-c.resp:
		w.Write(data.Value)
		clean()
		break
	case <-time.After(ki.timeout):
		w.WriteHeader(http.StatusRequestTimeout)
		io.WriteString(w, "Timeout")
		clean()
	}

}

func (ki *KafkaIntegrator) attachContext(input []byte, cid string) (attached bool, result []byte) {
	var dat map[string]interface{}
	attached = false
	result = input

	if err := json.Unmarshal(input, &dat); err == nil {

		hasContextWithCorrelation := util.MatchNestedMapPath(dat, correlationPath)
		hasContex := util.MatchNestedMapPath(dat, []string{CONTEXT})

		switch {
		case hasContextWithCorrelation:
			return
		case hasContex:
			context := dat[CONTEXT].(map[string]interface{})
			context[CORRELATION] = cid
			attached = true
			break
		case !hasContextWithCorrelation:
			context := make(map[string]interface{})
			context[CORRELATION] = cid
			attached = true

			dat[CONTEXT] = context
			break
		default:
			return

		}

		var enriched []byte
		if enriched, err = json.Marshal(dat); err != nil {
			panic(err)
		}
		result = enriched
	}

	return
}

func (ki *KafkaIntegrator) publishToKafka(data []byte, cid string) {
	fmt.Println("Publishing to Kafka...")
	fmt.Println(string(data))

	deliveryChan := make(chan kafka.Event)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &ki.cfg.PublishTopic,
			Partition: kafka.PartitionAny},
		Value:   data,
		Headers: []kafka.Header{kafka.Header{Key: CORRELATION, Value: []byte(cid)}},
	}

	ki.prod.Produce(msg, deliveryChan)

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
		fmt.Printf("[KAFKA] Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("[KAFKA] Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
}
