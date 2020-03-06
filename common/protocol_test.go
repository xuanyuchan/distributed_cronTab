package common

import (
	"bytes"
	"testing"
)

func TestJobUnMarshal(t *testing.T) {
	content := []byte(`{"name":"job1","command":"echo hello","cronExpr":"******"}`)
	reader := bytes.NewReader(content)
	job, err := UnMarshalJob(reader)
	if err != nil {
		t.Errorf("unmarshal error: %v", err)
	}
	if job.Name != "job1" {
		t.Errorf("unmarshal job name wrong, actual: %s", job.Name)
	}
	if job.Command != "echo hello" {
		t.Errorf("unmarshal job name wrong, actual: %s", job.Command)
	}
	if job.CronExpr != "******" {
		t.Errorf("unmarshal job name wrong, actual: %s", job.CronExpr)
	}
}
