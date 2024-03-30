package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

type LogTemplate struct {
	Message   string `json:"message"`
	Host      string `json:"host"`
	TimeStamp string `json:"@timestamp"`
	Level     string `json:"level"`
}

// 将日志信息保存到es中
func SaveLogToEs(message, host, timeStamp, level string) {
	//创建日志模板
	log2Save := LogTemplate{
		Message:   message,
		Host:      host,
		TimeStamp: timeStamp,
		Level:     level,
	}
	//创建es客户端
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Println("Failed to create es client:", err)
	}

	//创建一个上下文对象，上下文对象用于在程序的各个部分之间传递截止日期、取消信号以及其他请求范围的值
	ctx := context.Background()

	//序列化为json
	logJson, _ := json.Marshal(log2Save)

	//向es的某个索引写入数据
	response, err := client.Index().
		Index("logscan").
		BodyString(string(logJson)).
		Do(ctx)
	if err != nil {
		log.Println("Failed to save log to es:", err)
	}
	fmt.Println("log saved to es:", response)
}
