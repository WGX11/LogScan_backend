package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

func main() {
	//kafka broker地址
	brokers := []string{"10.134.150.136:9092"}

	//创建生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	fmt.Println("yes")
	if err != nil {
		log.Println("Failed to start Sarama producer:", err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Println("Failed to close Sarama producer:", err)
		}
	}()

	//构建并发送消息
	msg := &sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.StringEncoder("Hello Kafka from Go"),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Println("Failed to send message:", err)
	}
	log.Printf("Message is storted in topic(%s)/partition(%d)/offset(%d)\n", msg.Topic, partition, offset)

}
