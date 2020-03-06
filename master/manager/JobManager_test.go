package manager

import (
	"distributed_cronTab/common"
	"distributed_cronTab/master/config"
	"testing"
	"time"
)

func TestJobManager_SaveJob(t *testing.T) {
	_ = config.ParseConfig("../master.json")
	jM, err := InitJobManager()
	if err != nil {
		t.Errorf("init job manager error: %v\n", err)
		return
	}
	job := &common.Job{
		Name:     "Job_test",
		Command:  "echo hello",
		CronExpr: "*/2 * * * * * *",
	}
	err = jM.SaveJob(job)
	if err != nil {
		t.Errorf("Job manager save error: %v", err)
	}
}

func TestJobManager_DeleteJob(t *testing.T) {
	_ = config.ParseConfig("../master.json")
	jM, err := InitJobManager()
	if err != nil {
		t.Errorf("init job manager error: %v\n", err)
		return
	}
	name := "Job_test"
	err = jM.DeleteJob(name)
	if err != nil {
		t.Errorf("Job manager delete error: %v", err)
	}
}

func TestJobManager_ListJobs(t *testing.T) {
	_ = config.ParseConfig("../master.json")
	jM, err := InitJobManager()
	if err != nil {
		t.Errorf("init job manager error: %v\n", err)
		return
	}
	jobs, err := jM.ListJobs()
	if err != nil {
		t.Errorf("Job manager delete error: %v", err)
	}
	t.Logf("jobs count: %d\n", len(jobs))
	for _, job := range jobs {
		t.Logf("jobs:%v\n", job)
	}

}

func TestJobManager_KillJob(t *testing.T) {
	_ = config.ParseConfig("../master.json")
	jM, err := InitJobManager()
	if err != nil {
		t.Errorf("init job manager error: %v\n", err)
		return
	}
	name := "Job_test"
	err = jM.KillJob(name)
	if err != nil {
		t.Errorf("kill job error: %v\n", err)
		return
	}
	//test kill cron job
	count, err := jM.listKillJobs(name)
	if err != nil {
		t.Errorf("kill job error: %v\n", err)
		return
	}
	if count != 1 {
		t.Errorf("kill job test wrong, %s in etcd, count not 1, actual: %d\n", name, count)
		return
	}
	time.Sleep(time.Second * 3) //wait 3 seconds
	count, err = jM.listKillJobs(name)
	if err != nil {
		t.Errorf("kill job error: %v\n", err)
		return
	}
	if count != 0 {
		t.Errorf("kill job test wrong, %s in etcd, count not 0 after 1 second, actual: %d\n", name, count)
		return
	}
}
