package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type ScheduleEvent struct {
	EventType int
	Job       *Job
}

type JobSchedulePlan struct {
	Job      *Job
	NextTime time.Time
}

type JobExecuteInfo struct {
	JobPlan    *JobSchedulePlan
	PlanTime   time.Time
	ActualTime time.Time
}

type JobDoneResult struct {
	Job       *Job
	Error     error
	OutPut    []byte
	StartTime time.Time
	EndTime   time.Time
}

func UnMarshalJobFromReader(reader io.Reader) (*Job, error) {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("job reader read error: %v", err)
		return nil, err
	}
	return UnMarshalJob(content)
}

func UnMarshalJob(bytes []byte) (*Job, error) {
	job := &Job{}
	err := json.Unmarshal(bytes, job)
	if err != nil {
		log.Printf("job unmarshal error: %v", err)
		return nil, err
	}
	return job, nil
}

func ExtractJobName(key string) string {
	return strings.TrimPrefix(key, JobSaveDir)
}

func BuildScheduleEvent(eventType int, job *Job) *ScheduleEvent {
	return &ScheduleEvent{
		EventType: eventType,
		Job:       job,
	}
}

func BuildJobSchedulePlan(job *Job) (*JobSchedulePlan, error) {
	expr, err := cronexpr.Parse(job.CronExpr)
	if err != nil {
		return nil, err
	}
	nextTime := expr.Next(time.Now())
	return &JobSchedulePlan{
		Job:      job,
		NextTime: nextTime,
	}, nil
}

func BuildExecutedInfo(plan *JobSchedulePlan) *JobExecuteInfo {
	return &JobExecuteInfo{
		JobPlan:    plan,
		PlanTime:   plan.NextTime,
		ActualTime: time.Now(),
	}
}

func BuildJobDoneResult(job *Job, err error, result []byte, startTime, endTime time.Time) *JobDoneResult {
	return &JobDoneResult{
		Job:       job,
		Error:     err,
		OutPut:    result,
		StartTime: startTime,
		EndTime:   endTime,
	}
}
