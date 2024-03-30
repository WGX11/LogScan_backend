package pkg

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"time"
)

// 存储日志的入口
func KafkaConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	//创建kafka消费者
	consumer, err := sarama.NewConsumer([]string{"10.134.150.136:9092"}, config)
	if err != nil {
		log.Println("Failed to start Sarama consumer:", err)
	}
	defer func(consumer sarama.Consumer) {
		err := consumer.Close()
		if err != nil {

		}
	}(consumer)

	partitionConsumer, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
	if err != nil {
		log.Println("Failed to start Sarama consumer:", err)
	}

	defer func(partitionConsumer sarama.PartitionConsumer) {
		err := partitionConsumer.Close()
		if err != nil {
			log.Println("Failed to close Sarama consumer:", err)
		}
	}(partitionConsumer)

	//读取kafka中的消息
	location, err := time.LoadLocation("Asia/Shanghai") // 加载中国时区
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	for msg := range partitionConsumer.Messages() {
		currentTime := time.Now().In(location).Format(time.RFC3339)
		SaveLogToEs(string(msg.Key), "no host", currentTime, string(msg.Value))
		log.Printf("key is  %s\n", string(msg.Key))
		log.Printf("Message is %s\n", string(msg.Value))
	}
}
