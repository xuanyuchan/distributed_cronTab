package manager

import (
	"context"
	"distributed_cronTab/common"
	"distributed_cronTab/master/config"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
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

	G_JobManager = &JobManager{client: cli}
	return G_JobManager, nil
}

func (m *JobManager) SaveJob(job *common.Job) error {
	jobkey := common.JobSaveDir + job.Name
	jobVal, err := json.Marshal(job)
	if err != nil {
		log.Printf("Job Manager, save job, marshal error:%v\n", err)
		return err
	}

	kv := clientv3.NewKV(m.client)
	putResp, err := kv.Put(context.TODO(), jobkey, string(jobVal))
	if err != nil {
		log.Printf("job Manager, save job error: %v\n", err)
	}
	log.Println("job Manager, save job response revision: ", putResp.Header.Revision)
	return nil
}

func (m *JobManager) DeleteJob(name string) error {
	jobKey := common.JobSaveDir + name
	kv := clientv3.NewKV(m.client)

	_, err := kv.Delete(context.Background(), jobKey)
	if err != nil {
		log.Printf("job Manager, delete job error: %v\n", err)
		return err
	}
	return nil
}

func (m *JobManager) ListJobs() ([]*common.Job, error) {
	kv := clientv3.NewKV(m.client)
	dirKey := common.JobSaveDir
	getResp, err := kv.Get(context.Background(), dirKey, clientv3.WithPrefix())
	if err != nil {
		log.Printf("Job Manager, list jobs error: %v\n", err)
		return nil, err
	}

	jobs := make([]*common.Job, 0)
	for _, keyVal := range getResp.Kvs {
		tempJob := &common.Job{}
		err := json.Unmarshal(keyVal.Value, tempJob)
		if err != nil {
			log.Printf("job manager, list jobs, unmarshal job err: %v, jobresp: %s\n", err, keyVal.Value)
			continue
		}
		jobs = append(jobs, tempJob)
	}
	return jobs, nil
}

func (m *JobManager) KillJob(name string) error {
	jobKey := common.JobKillDir + name
	kv := clientv3.NewKV(m.client)
	lease := clientv3.NewLease(m.client)
	leaseResp, err := lease.Grant(context.Background(), 1)
	if err != nil {
		log.Printf("job manager, kill job error: %v\n", err)
		return err
	}
	_, err = kv.Put(context.Background(), jobKey, "", clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Printf("job manager, kill job error: %v\n", err)
		return err
	}
	return nil
}

func (m *JobManager) listKillJobs(name string) (int, error) {
	jobKey := common.JobKillDir + name
	kv := clientv3.NewKV(m.client)
	getResp, err := kv.Get(context.Background(), jobKey)
	if err != nil {
		log.Printf("job manager, list kill jobs error:%v\n", err)
		return 0, err
	}
	return int(getResp.Count), nil
}
