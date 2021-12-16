package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"s7i.io/kafka-gateway/internal/kafka"
)

func main() {

	prod := kafka.KafkaProducer{}
	prod.Init(&kafka.Properties{Server: "localhost:9092"})

	fmt.Println("running server")

	http.HandleFunc("/hello", hndHello)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}

func hndHello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello world @"+time.Now().String())

	fmt.Println(req)
}
