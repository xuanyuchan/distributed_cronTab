package executer

import (
	"distributed_cronTab/common"
	"testing"
)

func TestExecuter_ExecuteLoop(t *testing.T) {
	executer := InitExecuter()
	doneChan := executer.GetDoneJob()
	job := &common.Job{
		Name:     "job_test",
		Command:  "echo hello;sleep 1",
		CronExpr: "*/2 * * * * * *",
	}
	executer.AddNewJob(job)
	doneResult := <-doneChan
	t.Logf("done result: %+v", doneResult)
}
