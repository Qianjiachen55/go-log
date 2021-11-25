package main

import (
	"github.com/Qianjiachen55/go-log/kafka"
	"github.com/Qianjiachen55/go-log/run"
	"github.com/Qianjiachen55/go-log/tailfile"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

//日志收集
// filebeat
//woc wo de shu ru fa mei le
//get logs from a speacil dir, then send to kafka
//wo zhi neng xie pin ying huozhe english
type Config struct {
	KafKaConfig `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafKaConfig struct {
	Address string	`ini:"address""`
	Topic string `ini:"topic"`
	ChanSize int64 `ini:"chan_size"`
}
type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}



func main()  {
	var configObj = new(Config)
	//cfg,err := ini.Load("./conf/config.ini")
	//

	//kafkaAddr := cfg.Section("kafka").Key("address").String()
	logrus.SetLevel(logrus.TraceLevel)
	err :=ini.MapTo(configObj,"./conf/config.ini")
	if err != nil {
		logrus.Error("load config file err: ", err)
		return
	}
	//fmt.Printf("%#v\n",configObj)

	//1. connect to kafka
	err = kafka.InitKafka([]string{configObj.KafKaConfig.Address},configObj.KafKaConfig.ChanSize)
	if err != nil{
		logrus.Error("init kafka error: ",err)
		return
	}
	logrus.Debug("init kafka success!")

	err = tailfile.Init(configObj.CollectConfig.LogFilePath)
	if err != nil{
		logrus.Error("init tail file error: ",err)
		return
	}

	logrus.Debug("init tailFile success!")

	err = run.Run()
	if err != nil{
		logrus.Error("run failed err: ",err)
		return
	}

}