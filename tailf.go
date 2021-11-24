package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

func main() {
	filename := `./xx.log`
	config := tail.Config{
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: 2,
		},
		ReOpen:      true,
		MustExist:   false,
		Poll:        true,
		Pipe:        false,
		RateLimiter: nil,
		Follow:      true,
		MaxLineSize: 0,
		Logger:      nil,
	}

	tails, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	var msg * tail.Line
	var ok bool

	for {
		msg, ok = <-tails.Lines
		if !ok{
			fmt.Println("tail file close reopen,filename:%s",tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("msg: ",msg.Text)
	}

}
