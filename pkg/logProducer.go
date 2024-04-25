package pkg

import "log"

// 定义日志消息生产者和消费者
type LogProducer struct {
	Name      string
	Consumers map[int]chan string
}

func (lp *LogProducer) AddConsumer(id int, consumer chan string) {
	lp.Consumers[id] = consumer
}

func (lp *LogProducer) DeleteConsumer(id int) {
	consumer, ok := lp.Consumers[id]
	if !ok {
		return
	}
	close(consumer)
	delete(lp.Consumers, id)
}

var LogProducerInstance *LogProducer

// 向每个consumer发送日志消息
func SendLogMessage(message string) {
	consumers := LogProducerInstance.Consumers
	for _, consumer := range consumers {
		//非阻塞发送消息，防止某些消费者异常导致整体消息发送阻塞
		select {
		case consumer <- message:
		default:
			log.Printf("consumer is full, message: %s", message)
		}
	}
}

func init() {
	LogProducerInstance = &LogProducer{
		Name:      "logProducer",
		Consumers: make(map[int]chan string),
	}
}
