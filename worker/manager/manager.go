package manager

import (
	"context"
	"distributed_cronTab/common"
	"distributed_cronTab/worker/config"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"log"
	"time"
)

type JobManager struct {
	client *clientv3.Client
}

var (
	G_JobManager *JobManager
)

func InitJobManager() (*JobManager, error) {
	if G_JobManager != nil {
		return G_JobManager, nil
	}

	//init
	clientCfg := clientv3.Config{
		Endpoints:   config.G_Config.EtcdEndPoints,
		DialTimeout: time.Duration(config.G_Config.EtcdClientTimeOut) * time.Millisecond,
	}
	log.Printf("init job Manager, with etcd endpoints: %v\n", config.G_Config.EtcdEndPoints)
	cli, err := clientv3.New(clientCfg)
	if err != nil {
		log.Printf("etcd client build error: %v", err)
		return nil, err
	}
	//log.Printf("initOk\n")
	G_JobManager = &JobManager{client: cli}
	return G_JobManager, nil
}

func (m *JobManager) WatchJobs() (<-chan *common.ScheduleEvent, error) {
	kv := clientv3.NewKV(m.client)
	getResp, err := kv.Get(context.Background(), common.JobSaveDir, clientv3.WithPrefix())
	if err != nil {
		log.Printf("worker job manager, watch jobs, read from etcd error: %v\n", err)
		return nil, err
	}
	scheEventChan := make(chan *common.ScheduleEvent)
	go func() {
		kvalues := getResp.Kvs
		for _, kvpair := range kvalues {
			jobContent := kvpair.Value
			job, err := common.UnMarshalJob(jobContent)
			if err != nil {
				log.Printf("worker job manager, job watch unmarshal job error: %v, json content: %s\n", err, jobContent)
				continue
			}
			scheEvent := common.BuildScheduleEvent(common.SCHEDULE_PUT, job)
			scheEventChan <- scheEvent
		}

		//log.Printf("current task ok, %v", getResp.Kvs)

		watchRevision := getResp.Header.Revision
		watcher := clientv3.NewWatcher(m.client)
		watchChan := watcher.Watch(context.Background(), common.JobSaveDir, clientv3.WithRev(watchRevision+1), clientv3.WithPrefix())
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					job, err := common.UnMarshalJob(event.Kv.Value)
					if err != nil {
						log.Printf("worker job manager, watch jobs, err:%v, job json:%s\n", err, event.Kv.Value)
						continue
					}
					schEvent := common.BuildScheduleEvent(common.SCHEDULE_PUT, job)
					scheEventChan <- schEvent
				case mvccpb.DELETE:
					jobName := common.ExtractJobName(string(event.Kv.Key))
					job := &common.Job{Name: jobName}
					schEvent := common.BuildScheduleEvent(common.SCHEDULE_DELETE, job)
					scheEventChan <- schEvent
				}
			}
		}
	}()
	return scheEventChan, nil
}
