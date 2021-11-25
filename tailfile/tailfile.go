package tailfile

import (
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

var (
	TailObj *tail.Tail
)

func Init(filename string)(err error) {
	config := tail.Config{
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		ReOpen:    true,
		Follow:    true,
		Poll:      true,
		MustExist: false,
	}
	TailObj, err = tail.TailFile(filename,config)
	if err!=nil{
		logrus.Error("tailfile: init error :", err)
		return err
	}
	return nil
}
