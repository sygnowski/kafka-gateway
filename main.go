package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"s7i.io/kafka-gateway/internal/config"
	"s7i.io/kafka-gateway/internal/kafint"
)

func main() {

	conf := config.ReadConfig()
	kint := kafint.NewKafkaIntegrator(conf)

	fmt.Println("running server")

	http.HandleFunc("/hello", hndHello)
	http.HandleFunc(conf.App.Context, kint.Publish)

	bind := fmt.Sprintf("0.0.0.0:%s", strconv.Itoa(conf.App.Port))

	if err := http.ListenAndServe(bind, nil); err != nil {
		panic(err)
	}

}

func hndHello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello world @"+time.Now().String())

	fmt.Println(req)
}
