package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Qianjiachen55/go-log/common"
	"github.com/Qianjiachen55/go-log/tailfile"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"
	"time"
)

var client *clientv3.Client

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("connect to etcd error: %v", err)
		return err
	}

	return nil
}

// pull logConfig
func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd by key: %s failed, err: %v", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len:0 conf from etcd by key:%s", key)
	}
	ret := resp.Kvs[0]
	fmt.Println(string(ret.Value))
	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logrus.Errorf("json unmarshal failed,err:%v", err)
	}
	return collectEntryList, err
}

// watch etcd conf
func WatchConf(key string) {
	for{
		watchCh := client.Watch(context.Background(), key)

		for wResp := range watchCh {
			logrus.Infof("get new conf from etcd!")
			var newConf []common.CollectEntry
			for _, evt := range wResp.Events {
				fmt.Printf("type:%s key:%s value:%s \n", evt.Type, evt.Kv.Key, evt.Kv.Value)

				if evt.Type == clientv3.EventTypeDelete{
					logrus.Warningf("etcd delete the key!!!")
					tailfile.SendNewConf(newConf)
				}
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json unmarshall conf failed, err:%v", err)
					continue
				}
				tailfile.SendNewConf(newConf) // zu se
			}
		}
	}

}


