package logic

import (
	"encoding/json"
	"gopkg.in/Shopify/sarama.v1"
	"log"
	"os"
	"robTickets/src/dao/mysql"
	"robTickets/src/models"
)

var (
	client      sarama.Client
	producerMap map[string]sarama.AsyncProducer
	logger      = log.New(os.Stderr, "", log.LstdFlags)
	//bufferSize  = flag.Int("buffer-size", 256, "The buffer size of the message channel.")
)

const (
	insertOrderTopic        string = "insertOrder"
	updateTicketNumberTopic string = "updateTicket"
	cancelOrderTopic        string = "cancelOrder"
	incr                    int8   = 1
	decr                    int8   = -1
)

func ProducerClose() {
	if err := producerMap[insertOrderTopic].Close(); err != nil {
		logger.Println("Failed to close Kafka producer insertOrder cleanly:", err)
	}
	if err := producerMap[updateTicketNumberTopic].Close(); err != nil {
		logger.Println("Failed to close Kafka producer updateTicket cleanly:", err)
	}
	if err := producerMap[cancelOrderTopic].Close(); err != nil {
		logger.Println("Failed to close Kafka producer cancelOrder cleanly:", err)
	}
	if err := client.Close(); err != nil {
		logger.Println("Failed to close Kafka client cleanly:", err)
	}
}

func InitKafka() error {
	config := sarama.NewConfig()
	//config.Producer.RequiredAcks = sarama.WaitForAll // Producer生产者: 发送完数据需要leader和follow都确认
	var err error
	// 连接Kafka
	client, err = sarama.NewClient([]string{"47.102.42.113:9092"}, config)
	if err != nil {
		return err
	}
	producerMap = make(map[string]sarama.AsyncProducer)
	producerMap[insertOrderTopic], err = sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return err
	}
	producerMap[updateTicketNumberTopic], err = sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return err
	}
	producerMap[cancelOrderTopic], err = sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return err
	}

	if err = consumerInsertOrder(); err != nil {
		return err
	}
	if err = consumeTicketNumber(); err != nil {
		return err
	}
	if err = consumeCancelOrder(); err != nil {
		return err
	}
	return nil
}

func SendMessage(topicName string, msg sarama.Encoder) {
	message := &sarama.ProducerMessage{
		Topic: topicName,
		Value: msg,
	}
	select {
	case producerMap[topicName].Input() <- message:
	case err := <-producerMap[topicName].Errors():
		logger.Println("Kafka发送信息失败: " + err.Error())
	}
}

func consumerInsertOrder() error {
	c, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return err
	}
	partitionList, err := c.Partitions(insertOrderTopic)

	//var (
	//	//messages = make(chan map[string]*sarama.ConsumerMessage, *bufferSize)
	//	wg sync.WaitGroup
	//)

	for _, partition := range partitionList {
		pc, err := c.ConsumePartition(insertOrderTopic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}
		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()

			for message := range pc.Messages() {
				//messages <- message
				var orderData models.TicketOrder
				err := json.Unmarshal(message.Value, &orderData)
				if err != nil {
					logger.Println("json反序列化失败!", err)
					continue
				}
				if err := mysql.CreateOrder(orderData); err != nil {
					logger.Println("mysql插入订单失败!", err)
				}
			}
		}(pc)
	}
	return nil
}

func consumeTicketNumber() error {
	c, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return err
	}
	partitionList, err := c.Partitions(updateTicketNumberTopic)

	for _, partition := range partitionList {
		pc, err := c.ConsumePartition(updateTicketNumberTopic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}
		//wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for message := range pc.Messages() {
				//messages <- message
				var ticketData models.MTicket
				err := json.Unmarshal(message.Value, &ticketData)
				if err != nil {
					logger.Println("json反序列化失败!", err)
					continue
				}
				if ticketData.Type == decr {
					if err := mysql.DescTicketNumber(ticketData.TicketID, ticketData.Quantity); err != nil {
						logger.Println("mysql减少票数失败!", err)
					}
				} else if ticketData.Type == incr {
					if err := mysql.IncrTicketNumber(ticketData.TicketID, ticketData.Quantity); err != nil {
						logger.Println("mysql增加票数失败!", err)
					}
				}
			}
		}(pc)
	}
	return nil
}

func consumeCancelOrder() error {
	c, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return err
	}
	partitionList, err := c.Partitions(cancelOrderTopic)

	for _, partition := range partitionList {
		pc, err := c.ConsumePartition(cancelOrderTopic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for message := range pc.Messages() {
				var cancel models.MCancelOrder
				err := json.Unmarshal(message.Value, &cancel)
				if err != nil {
					logger.Println("json反序列化失败!", err)
					continue
				}
				if err := mysql.CancelOrder(cancel.OrderID); err != nil {
					logger.Println("mysql取消票数失败!", err)
				}
			}
		}(pc)
	}
	return nil
}
