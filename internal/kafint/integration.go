package kafint

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"s7i.io/kafka-gateway/internal/config"
	"s7i.io/kafka-gateway/internal/util"
)

const CORRELATION string = "correlation"
const CONTEXT string = "context"

func CorrelationPath() []string {
	return []string{CONTEXT, CORRELATION}
}

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
	cfg            *config.Cfg
	correlationMap sync.Map
	timeout        time.Duration
}

func NewKafkaIntegrator(config *config.Cfg) *KafkaIntegrator {
	kint := KafkaIntegrator{}
	kint.init(config)
	return &kint
}

func (ki *KafkaIntegrator) init(config *config.Cfg) {
	ki.cfg = config
	ki.timeout = time.Duration(config.App.Timeout) * time.Second
	fmt.Printf("Timeout :%s\n", ki.timeout)

	prodProps := &kafka.ConfigMap{}
	append(&config.Pub, prodProps)

	consProps := &kafka.ConfigMap{}
	append(&config.Sub, consProps)

	fmt.Println("making kafka producer")

	p, err := kafka.NewProducer(prodProps)

	if err != nil {

		panic(err)
	}
	ki.prod = p

	c, err := kafka.NewConsumer(consProps)
	if err != nil {
		panic(err)
	}
	ki.cons = c

	ki.correlationMap = sync.Map{}

	c.Subscribe(config.Sub.Topic, nil)
	go ki.fetch()
}

func append(spec *config.CfgSpec, configMap *kafka.ConfigMap) {
	for _, p := range spec.Properties {
		configMap.Set(p)
	}
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
					println("[KI] got cid in the body.context.correlation: " + cid.(string))
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
				println("[KI] got cid in the header: " + cid)

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
	var ctxAtt, cidExists bool
	var cid string

	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer req.Body.Close()
	}

	if cidExists, cid = correlationInBody(bodyBytes); !cidExists && cid == "" {
		cid = util.UUID()
		ctxAtt, bodyBytes = ki.attachContext(bodyBytes, cid)
		if ctxAtt {
			fmt.Printf("[KI] New context attached [%s].\n", cid)
		}
	}
	ki.publishToKafka(bodyBytes, cid)

	c := &Correlaction{id: cid, resp: make(chan kafka.Message, 1)}
	ki.correlationMap.Store(cid, c)

	clean := func() {
		close(c.resp)
		ki.correlationMap.Delete(cid)
	}
	defer clean()

	select {
	case data := <-c.resp:
		w.Write(data.Value)
		break
	case <-time.After(ki.timeout):
		statusGatewayTimeout(w)
	}
}

func statusGatewayTimeout(w http.ResponseWriter) {
	w.WriteHeader(http.StatusGatewayTimeout)
	io.WriteString(w, "Gatewat Timeout.")
}

func (ki *KafkaIntegrator) attachContext(input []byte, cid string) (attached bool, result []byte) {
	var dat map[string]interface{}
	attached = false
	result = input

	if err := json.Unmarshal(input, &dat); err == nil {

		hasContextWithCorrelation := util.MatchNestedMapPath(dat, CorrelationPath())
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
	fmt.Printf("[KI] Publishing: [%s], correlation[%s]\n", string(data), cid)

	deliveryChan := make(chan kafka.Event)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &ki.cfg.Pub.Topic,
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

func correlationInBody(body []byte) (bool, string) {
	jsonMap := make(map[string]interface{})

	if err := json.Unmarshal(body, &jsonMap); err == nil {
		full, ctx := util.GetLastMapPath(jsonMap, CorrelationPath())

		if full {
			cid := ctx.(string)
			fmt.Printf("[KI] Extracted CID [%s], full-path [%v].\n", cid, full)
			return true, cid
		}

	} else {
		panic(err)
	}
	return false, ""
}
