package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"s7i.io/kafka-gateway/internal/kafint"
)

func main() {

	prod := kafint.KafkaIntegrator{}
	prod.Init(&kafint.Properties{
		Server:           os.Getenv("BROKER"),
		PublishTopic:     os.Getenv("PUBLISH_TOPIC"),
		SubscribeTopic:   os.Getenv("SUBSCRIBE_TOPIC"),
		SubscribeGroupId: os.Getenv("SUBSCRIBE_GROUP_ID"),
	})

	fmt.Println("running server")

	http.HandleFunc("/hello", hndHello)
	http.HandleFunc("/", prod.Publish)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}

func hndHello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello world @"+time.Now().String())

	fmt.Println(req)
}
