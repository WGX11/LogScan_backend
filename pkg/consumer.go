package pkg

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// 存储日志的入口
func KafkaConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	//创建kafka消费者
	//consumer, err := sarama.NewConsumer([]string{"10.134.150.136:9092"}, config)
	consumer, err := sarama.NewConsumer([]string{"10.8.0.6:9092"}, config)
	if err != nil {
		log.Println("Failed to start Sarama consumer:", err)
		return
	}
	defer func(consumer sarama.Consumer) {
		if consumer == nil {
			log.Println("consumer is nil")
			return
		}
		err := consumer.Close()
		if err != nil {
			log.Println("Failed to close consumer:", err)
			return
		}
	}(consumer)

	partitionConsumer, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
	if err != nil {
		log.Println("Failed to start Sarama consumer:", err)
		return
	}

	defer func(partitionConsumer sarama.PartitionConsumer) {
		err := partitionConsumer.Close()
		if err != nil {
			log.Println("Failed to close Sarama consumer:", err)
			return
		}
	}(partitionConsumer)

	//读取kafka中的消息
	location, err := time.LoadLocation("Asia/Shanghai") // 加载中国时区
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	for msg := range partitionConsumer.Messages() {
		//向消费者发送消息
		SendLogMessage(string(msg.Key))
		currentTime := time.Now().In(location).Format(time.RFC3339)
		randomNumber := rand.Intn(50)
		randomNumberString := strconv.Itoa(randomNumber)
		SaveLogToEs(string(msg.Key), "host"+randomNumberString, currentTime, string(msg.Value))
		log.Printf("key is  %s\n", string(msg.Key))
		log.Printf("Message is %s\n", string(msg.Value))
	}
}
