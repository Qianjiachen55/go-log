package tailfile

import (
	"github.com/Qianjiachen55/go-log/kafka"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"time"
)

type tailTask struct {
	path  string
	topic string
	tObj  *tail.Tail
}



func newTailTask(path, topic string) *tailTask {
	tt := &tailTask{
		path:  path,
		topic: topic,
	}
	return tt
}
func (t *tailTask) Init() (err error) {
	config := tail.Config{
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		ReOpen:    true,
		Follow:    true,
		Poll:      true,
		MustExist: false,
	}

	t.tObj, err = tail.TailFile(t.path, config)
	return err
}

func (t *tailTask) run() {
	logrus.Infof("collect for path %s is running! ",t.path)
	for {
		line, ok := <-t.tObj.Lines
		if !ok {
			logrus.Warn("tail file clost reopen, filename: ", t.path)
			time.Sleep(time.Second)
			continue
		}
		//
		msg := &sarama.ProducerMessage{}
		msg.Topic = t.topic
		msg.Value = sarama.StringEncoder(line.Text)

		kafka.ToMsgChan(msg)
	}
}


