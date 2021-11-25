package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	MsgChan chan *sarama.ProducerMessage
)

func InitKafka(address []string, chanSize int64) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer(address, config)

	if err != nil {
		logrus.Error("Kafka:producer closed err: ", err)
		return err
	}
	MsgChan = make(chan *sarama.ProducerMessage, chanSize)
	// go sendMsg hei hei
	go sendMsg()
	return nil
}

// read from msgChan then send to kafka
func sendMsg() {
	for {
		select {
		case msg := <-MsgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Warn("send to kafka failed ,err: ", err)
				return
			}
			logrus.Debugf("send msg to kafka . pid :%v ,offset:%v", pid, offset)
		}
	}
}
