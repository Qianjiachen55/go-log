package main

import (
	"fmt"
	"github.com/Qianjiachen55/go-log/etcd"
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
	EtcdConfig `ini:"etcd"`
}

type KafKaConfig struct {
	Address string	`ini:"address"`
	ChanSize int64 `ini:"chan_size"`
}
type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

type EtcdConfig struct {
	Address string `ini:"address"`
	CollectKey string `ini:"collect_key"`
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
	// init etcd
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil{
		logrus.Error("init etcd error: ",err)
		return
	}

	// pull  config from etcd

	allConf,err :=etcd.GetConf(configObj.EtcdConfig.CollectKey)
	if err != nil{
		logrus.Errorf("get conf from etcd failed: err : %v",err)
	}
	fmt.Println(allConf)

	//watch etcd

	go etcd.WatchConf(configObj.EtcdConfig.CollectKey)
	err = tailfile.Init(allConf)
	if err != nil{
		logrus.Error("init tail file error: ",err)
		return
	}

	logrus.Debug("init tailFile success!")

	run.Run()

}