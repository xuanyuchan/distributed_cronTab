package scheduler

import (
	"distributed_cronTab/common"
	"log"
	"testing"
	"time"
)

func TestScheduler_ScheduleLoop(t *testing.T) {
	eventArr := []*common.ScheduleEvent{
		&common.ScheduleEvent{
			EventType: common.SCHEDULE_PUT,
			Job: &common.Job{
				Name:     "job_test_1",
				Command:  "echo hello",
				CronExpr: "*/2 * * * * * *",
			},
		},
		&common.ScheduleEvent{
			EventType: common.SCHEDULE_PUT,
			Job: &common.Job{
				Name:     "job_test_1",
				Command:  "echo hello 1",
				CronExpr: "*/2 * * * * * *",
			},
		},
		&common.ScheduleEvent{
			EventType: common.SCHEDULE_DELETE,
			Job: &common.Job{
				Name:     "job_test_1",
				Command:  "echo hello",
				CronExpr: "*/2 * * * * * *",
			},
		},
	}

	s, err := InitScheduler()
	if err != nil {
		t.Errorf("init scheduler error:%v\n", err)
		return
	}

	//idx := 0
	go func() {
		ch := s.GetJobNeedScheduled()
		for job := range ch {
			log.Printf("schedule Job: %+v", job)
		}
	}()
	//for _, e := range eventArr {
	//	s.Submit(e)
	//	time.Sleep(time.Second * 1)
	//	t.Logf("plan table :%v\n", s.scheduleTable)
	//	idx++
	//	if idx == 1 {
	//		if len(s.scheduleTable) != 1 {
	//			t.Errorf("schedule loop test wrong, actual count: %d\n", len(s.scheduleTable))
	//		}
	//	} else if idx == 2 {
	//		if jobPlan, ok := s.scheduleTable["job_test_1"]; !ok || jobPlan.Job.Command != "echo hello 1" {
	//			t.Errorf("scheduler loop test wrong, modify error, actual jobPlan: %+v\n", jobPlan)
	//		}
	//	} else if idx == 3 {
	//		if len(s.scheduleTable) != 0 {
	//			t.Errorf("schedule loop test wrong, delete wrong count: %d\n", len(s.scheduleTable))
	//		}
	//	}
	//}
	s.Submit(eventArr[1])
	go func() {
		for {
			//log.Printf("schedule map: %+v", s.scheduleTable)
			s.FinishJob(eventArr[2].Job)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	time.Sleep(100 * time.Second)
}

func TestTryScheduler(t *testing.T) {
	event := &common.ScheduleEvent{
		EventType: common.SCHEDULE_PUT,
		Job: &common.Job{
			Name:     "job_test_1",
			Command:  "echo hello",
			CronExpr: "*/5 * * * * * *",
		},
	}
	s, err := InitScheduler()
	if err != nil {
		t.Errorf("init scheduler error:%v\n", err)
		return
	}
	//s.StartSchedule()
	s.Submit(event)
	for {
		time.Sleep(500 * time.Millisecond)
		jobPlan, nextDuration := s.trySchedule()
		log.Printf("job: %+v, next time: %v\n", jobPlan, time.Now().Add(nextDuration))
	}
}
