package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
	"time"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true


	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("producer closed,err: ", err)
		return
	}
	defer client.Close()

	msg := &sarama.ProducerMessage{}
	msg.Topic = "shopping"
	msg.Value = sarama.StringEncoder("this is a test log"+cast.ToString(time.Now()))

	message, offset, err := client.SendMessage(msg)

	if err !=nil{
		fmt.Println("send failed: ",err)
		return
	}
	fmt.Println("res: ",message," ",offset)

}
