package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"logscan/internal/handler"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/chartData", handler.GetBarChartData)
	router.GET("/searchData", handler.LogMessageInfoHandler)
	err := router.Run(":9031")
	if err != nil {
		return
	}
}
