package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"s7i.io/kafka-gateway/internal/kafka"
)

func main() {

	prod := kafka.KafkaProducer{}
	prod.Init(&kafka.Properties{
		Server:       os.Getenv("BROKER"),
		PublishTopic: os.Getenv("PUBLISH_TOPIC"),
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
