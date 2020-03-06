package main

import (
	"distributed_cronTab/common"
	"distributed_cronTab/worker/config"
	"distributed_cronTab/worker/executer"
	"distributed_cronTab/worker/manager"
	"distributed_cronTab/worker/scheduler"
	"fmt"
	"log"
)

type JobManager interface {
	WatchJobs() (<-chan *common.ScheduleEvent, error)
}

type Scheduler interface {
	Submit(*common.ScheduleEvent)
	FinishJob(job *common.Job)
	GetJobNeedScheduled() chan *common.Job
}

type Executer interface {
	AddNewJob(*common.Job)
	GetDoneJob() chan *common.JobDoneResult
}

func initConfig() {
	err := config.ParseConfig("./worker/worker.json")
	if err != nil {
		panic(err)
	}
}

func initJobManager() (JobManager, error) {
	return manager.InitJobManager()
}

func initScheduler() (Scheduler, error) {
	return scheduler.InitScheduler()
}

func initExecuter() Executer {
	return executer.InitExecuter()
}

func main() {
	fmt.Println("worker start")
	initConfig()
	jM, err := initJobManager()
	if err != nil {
		panic(err)
	}
	sch, err := initScheduler()
	if err != nil {
		panic(err)
	}
	ch, err := jM.WatchJobs()
	if err != nil {
		panic(err)
	}
	executer := initExecuter()

	executerChan := sch.GetJobNeedScheduled()
	jobDoneChan := executer.GetDoneJob()

	for {
		select {
		case event := <-ch:
			sch.Submit(event)
		case executeJob := <-executerChan:
			executer.AddNewJob(executeJob)
		case jobDoneResult := <-jobDoneChan:
			log.Printf("Job done result: %+v\n", jobDoneResult)
			sch.FinishJob(jobDoneResult.Job)
		}
	}

}
