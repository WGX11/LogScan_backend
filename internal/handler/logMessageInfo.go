package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"logscan/pkg"
)

type LogMessage struct {
	Message   string `json:"message"`
	Host      string `json:"host"`
	TimeStamp string `json:"@timestamp"`
	Level     string `json:"level"`
}

// 搜索页面日志信息
func LogMessageInfoHandler(ctx *gin.Context) {

	// 预检请求的处理
	startTime := ctx.Query("start")
	endTime := ctx.Query("end")
	lucene := ctx.Query("lucene")
	response := pkg.SearchLogFromEs(startTime, endTime, lucene)
	logMessages := make([]LogMessage, 0)
	for _, hit := range response.Hits.Hits {
		logEntry := make(map[string]interface{})
		if err := json.Unmarshal(hit.Source, &logEntry); err != nil {
			log.Println("Failed to unmarshal hit source:", err)
			continue
		}
		var timeStamp string
		timeStamp1, ok1 := logEntry["@timestamp"].(string)
		timeStamp2, ok2 := logEntry["timestamp"].(string)
		if ok1 {
			timeStamp = timeStamp1
		} else if ok2 {
			timeStamp = timeStamp2
		}

		var message string
		message1, ok1 := logEntry["message"].(string)
		message2, ok2 := logEntry["msg"].(string)
		if ok1 {
			message = message1
		} else if ok2 {
			message = message2
		}

		var host string
		host1, ok1 := logEntry["host"].(string)
		host2, ok2 := logEntry["hostname"].(string)
		if ok1 {
			host = host1
		} else if ok2 {
			host = host2
		}

		var level string
		level1, ok1 := logEntry["level"].(string)
		level2, ok2 := logEntry["severity"].(string)
		if ok1 {
			level = level1
		} else if ok2 {
			level = level2
		}
		logMessages = append(logMessages, LogMessage{
			Message:   message,
			Host:      host,
			TimeStamp: timeStamp,
			Level:     level,
		})
	}
	ctx.JSON(200, logMessages)
}
