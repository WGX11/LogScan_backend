package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"logscan/internal/handler"
	"logscan/pkg"
	"logscan/sql"
)

func main() {
	//初始化数据库
	sql.InitDB()
	defer sql.CloseDB()
	//启动日志接收
	go func() {
		pkg.KafkaConsumer()
	}()
	//启动报警监控
	go func() {
		pkg.StartAlarmMonitor()
	}()
	pkg.TestMonitor()
	fmt.Println("after consumer")
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/chartData", handler.GetBarChartData)
	router.GET("/searchData", handler.LogMessageInfoHandler)
	router.GET("/anomalyData", handler.GetNotificationData)
	router.GET("/monitorData", handler.GetMonitorData)
	router.GET("/alarmList", handler.HandleGetAlarmList)
	router.GET("/alarmDashBoard", handler.HandleGetAlarmDashBoard)

	router.POST("/alarmAdd", handler.HandleAddAlarm)

	router.DELETE("/alarmDelete/:id", handler.HandleDeleteAlarm)
	err := router.Run(":9031")
	if err != nil {
		return
	}
}
