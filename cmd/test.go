package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"logscan/internal/handler"
	"logscan/pkg"
)

func main() {
	go func() {
		pkg.KafkaConsumer()
	}()
	fmt.Println("after consumer")
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/chartData", handler.GetBarChartData)
	router.GET("/searchData", handler.LogMessageInfoHandler)
	router.GET("/anomalyData", handler.GetNotificationData)
	router.GET("/monitorData", handler.GetMonitorData)
	err := router.Run(":9031")
	if err != nil {
		return
	}
}
