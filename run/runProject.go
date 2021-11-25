package run

import (
	"github.com/Qianjiachen55/go-log/kafka"
	"github.com/Qianjiachen55/go-log/tailfile"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"time"
)

func Run() (err error) {
	for {
		line, ok := <- tailfile.TailObj.Lines
		if !ok{
			logrus.Warn("tail file clost reopen, filename: ",tailfile.TailObj.Filename)
			time.Sleep(time.Second)
			continue
		}
		//
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web"
		msg.Value = sarama.StringEncoder(line.Text)

		kafka.MsgChan <- msg
	}

	return nil
}
