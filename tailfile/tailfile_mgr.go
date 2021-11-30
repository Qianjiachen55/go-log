package tailfile

import (
	"github.com/Qianjiachen55/go-log/common"
	"github.com/sirupsen/logrus"
)

type tailTaskMgr struct {
	tailTaskMap map[string]*tailTask
	collectEntryList []common.CollectEntry
	confChan chan []common.CollectEntry
}

var (
	ttMgr *tailTaskMgr
)


func Init(allConf []common.CollectEntry) (err error) {
	ttMgr = &tailTaskMgr{
		tailTaskMap:  	  map[string]*tailTask{},
		collectEntryList: allConf,
		confChan:         make(chan []common.CollectEntry),
	}

	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic)
		err = tt.Init()
		if err != nil {
			logrus.Errorf("create tailObj for path:%s failed,err:%v", conf.Path, err)
			continue
		}
		// shou ji ri zhi
		logrus.Infof("create a tail task for path:%s success", conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt
		go tt.run()
	}
	go ttMgr.watch() //waiting new config

	return nil
}

func (t *tailTaskMgr) watch() {
	for{
		newConf := <-t.confChan
		logrus.Infof("get new conf from etcd, conf:%v", newConf)
		for _, conf := range newConf {
			if t.isExist(conf){
				continue
			}
			//---//
			tt := newTailTask(conf.Path, conf.Topic)
			err := tt.Init()
			if err != nil {
				logrus.Errorf("create tailObj for path:%s failed,err:%v", conf.Path, err)
				continue
			}
			// shou ji ri zhi
			logrus.Infof("create a tail task for path:%s success", conf.Path)
			ttMgr.tailTaskMap[tt.path] = tt

			go tt.run()
			//---//

		}
	}

}

func (t *tailTaskMgr)isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}

func SendNewConf(newConf []common.CollectEntry){
	ttMgr.confChan <- newConf
}