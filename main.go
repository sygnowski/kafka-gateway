package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"s7i.io/kafka-gateway/internal/config"
	"s7i.io/kafka-gateway/internal/kafint"
)

func main() {

	conf := config.ReadConfig()
	kint := kafint.NewKafkaIntegrator(conf)

	eng := gin.Default()

	bind := fmt.Sprintf("0.0.0.0:%s", strconv.Itoa(conf.App.Port))

	eng.POST(conf.App.Context, func(ctx *gin.Context) {
		kint.Publish(ctx.Writer, ctx.Request)
	})
	eng.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world @"+time.Now().String())
	})
	eng.Run(bind)

}
