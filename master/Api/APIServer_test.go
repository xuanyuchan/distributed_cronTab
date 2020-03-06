package Api

import (
	"bytes"
	"distributed_cronTab/common"
	"distributed_cronTab/master/config"
	"distributed_cronTab/master/manager"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveJobFunc(t *testing.T) {
	_ = config.ParseConfig("../master.json")
	j, _ := manager.InitJobManager()
	InitApiServer(j)
	srv := httptest.NewServer(http.HandlerFunc(jobSaveHandler))
	defer srv.Close()

	//param = `{"name":"jobTest","command":"echo hello","cronExpr":"*****"}`
	job := &common.Job{
		Name:     "job_test",
		Command:  "echo hello",
		CronExpr: "*****",
	}
	param, _ := json.Marshal(job)
	resp, err := http.Post(srv.URL, "application/json", bytes.NewBuffer(param))
	defer resp.Body.Close()
	if err != nil {
		t.Errorf("testAPIServer save job, error:%v\n", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	t.Logf("resp: %s\n", body)
	apiResponse := &apiResponse{}
	json.Unmarshal(body, apiResponse)
	if apiResponse.ErrNo != 0 {
		t.Errorf("testAPIServer save job, errorno: %d\n, error msg:%s\n", apiResponse.ErrNo, apiResponse.ErrMsg)
	}
}

func TestDeleteJob(t *testing.T) {
	_ = config.ParseConfig("../master.json")
	j, _ := manager.InitJobManager()
	InitApiServer(j)
	srv := httptest.NewServer(http.HandlerFunc(JobDeleteHandler))
	defer srv.Close()

	reqContent := []byte(`{"name":"job_test"}`)
	resp, err := http.Post(srv.URL, "application/json", bytes.NewBuffer(reqContent))
	if err != nil {
		t.Errorf("testAPIServer delete job, error:%v\n", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	t.Logf("resp: %s\n", body)
	apiResponse := &apiResponse{}
	json.Unmarshal(body, apiResponse)
	if apiResponse.ErrNo != 0 {
		t.Errorf("testAPIServer delete job, errorno: %d\n, error msg:%s\n", apiResponse.ErrNo, apiResponse.ErrMsg)
	}
}
