package executer

import (
	"context"
	"distributed_cronTab/common"
	"distributed_cronTab/common/util"
	"log"
	"os/exec"
	"time"
)

type Executer struct {
	jobCh     chan *common.Job
	jobDoneCh chan *common.JobDoneResult
}

var (
	G_Executer *Executer
)

func (e *Executer) ExecuteLoop() {
	for {
		select {
		case job := <-e.jobCh:
			go doExecuteJob(job, e.jobDoneCh)
		}
	}
}

func (e *Executer) AddNewJob(job *common.Job) {
	go func() { e.jobCh <- job }()
}

func doExecuteJob(job *common.Job, done chan *common.JobDoneResult) {
	//lock
	lock, err := util.InitLock(job.Name)
	if err != nil {
		log.Printf("job executer lock error: %v\n", err)
		jobResult := common.BuildJobDoneResult(job, EXECUTE_LOCK_ERR, nil, time.Now(), time.Now())
		done <- jobResult
		return
	}
	err = lock.Lock()
	defer lock.UnLock()
	if err != nil {
		log.Printf("job executer lock error: %v\n", err)
		jobResult := common.BuildJobDoneResult(job, EXECUTE_LOCK_ERR, nil, time.Now(), time.Now())
		done <- jobResult
		return
	}
	command := job.Command
	startTime := time.Now()
	cmd := exec.CommandContext(context.Background(), "/bin/bash", "-c", command)
	result, err := cmd.Output()
	endTime := time.Now()
	jobResult := common.BuildJobDoneResult(job, err, result, startTime, endTime)
	done <- jobResult
}

func (e *Executer) GetDoneJob() chan *common.JobDoneResult {
	return e.jobDoneCh
}

func InitExecuter() *Executer {
	if G_Executer != nil {
		return G_Executer
	}
	G_Executer = &Executer{
		jobCh:     make(chan *common.Job),
		jobDoneCh: make(chan *common.JobDoneResult),
	}
	go G_Executer.ExecuteLoop()
	return G_Executer
}
