package scheduler

import (
	"distributed_cronTab/common"
	"log"
	"time"
)

type Scheduler struct {
	ch               chan *common.ScheduleEvent
	jobNeedScheduled chan *common.Job
	scheduleTable    map[string]*common.JobSchedulePlan
	executedTable    map[string]*common.JobExecuteInfo
}

var (
	G_Scheduler *Scheduler
)

func (s *Scheduler) Submit(event *common.ScheduleEvent) {
	s.ch <- event
}

func (s *Scheduler) GetJobNeedScheduled() chan *common.Job {
	return s.jobNeedScheduled
}

func (s *Scheduler) ScheduleLoop() {
	jobPlan, nextDuration := s.trySchedule()
	for {
		select {
		case event := <-s.ch:
			s.handleScheduleEvent(event)
		case <-time.After(nextDuration):
		}
		jobPlan, nextDuration = s.trySchedule()
		if jobPlan != nil {
			//log.Printf("start doing task: %+v", jobPlan.Job)
			job := s.tryExecute(jobPlan)
			if job != nil {
				s.jobNeedScheduled <- job
			}
		}
	}
}

func (s *Scheduler) trySchedule() (*common.JobSchedulePlan, time.Duration) {
	now := time.Now()
	nextDuration := time.Duration(0)
	for jobName, jobPlan := range s.scheduleTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			job := jobPlan.Job
			newPlan, _ := common.BuildJobSchedulePlan(job)
			s.scheduleTable[jobName] = newPlan
			return jobPlan, 0
		}
		if nextDuration == 0 || jobPlan.NextTime.Sub(now) < time.Duration(nextDuration) {
			nextDuration = jobPlan.NextTime.Sub(now)
		}
	}
	return nil, nextDuration
}

func (s *Scheduler) tryExecute(jobPlan *common.JobSchedulePlan) *common.Job {
	if _, executing := s.executedTable[jobPlan.Job.Name]; executing {
		log.Printf("job: %s running, skipped executing\n", jobPlan.Job.Name)
		return nil
	}
	executeInfo := common.BuildExecutedInfo(jobPlan)
	s.executedTable[jobPlan.Job.Name] = executeInfo
	return jobPlan.Job
}

func (s *Scheduler) FinishJob(job *common.Job) {
	delete(s.executedTable, job.Name)
}

func (s *Scheduler) handleScheduleEvent(event *common.ScheduleEvent) {
	switch event.EventType {
	case common.SCHEDULE_PUT:
		jobPlan, err := common.BuildJobSchedulePlan(event.Job)
		if err != nil {
			log.Printf("scheduler build job plan, job:+%v, err: %v\n", event.Job, err)
			return
		}
		s.scheduleTable[event.Job.Name] = jobPlan
	case common.SCHEDULE_DELETE:
		if _, ok := s.scheduleTable[event.Job.Name]; ok {
			delete(s.scheduleTable, event.Job.Name)
		}
	}
}

func InitScheduler() (*Scheduler, error) {
	if G_Scheduler != nil {
		return G_Scheduler, nil
	}
	channel := make(chan *common.ScheduleEvent)
	G_Scheduler = &Scheduler{
		ch:               channel,
		scheduleTable:    make(map[string]*common.JobSchedulePlan),
		jobNeedScheduled: make(chan *common.Job),
		executedTable:    make(map[string]*common.JobExecuteInfo),
	}
	go G_Scheduler.ScheduleLoop()
	return G_Scheduler, nil
}

//func (s *Scheduler) StartSchedule() {
//	go s.ScheduleLoop()
//}

func (s *Scheduler) stopScheduler() {
	close(s.ch)
}
