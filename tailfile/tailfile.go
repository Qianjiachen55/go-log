package tailfile

import (
	"context"
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
	ctx context.Context
	cancel context.CancelFunc

}



func newTailTask(path, topic string) *tailTask {
	ctx,cancel := context.WithCancel(context.Background())
	tt := &tailTask{
		path:  path,
		topic: topic,
		ctx: ctx,
		cancel: cancel,
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
		select {
		case <- t.ctx.Done():
			return
		case line, ok := <-t.tObj.Lines:
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
}


