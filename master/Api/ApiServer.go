package Api

import (
	"distributed_cronTab/common"
	"distributed_cronTab/master/config"
	"net"
	"net/http"
	"strconv"
	"time"
)

type JobOperator interface {
	SaveJob(job *common.Job) error
	DeleteJob(name string) error
	ListJobs() ([]*common.Job, error)
	KillJob(name string) error
}

type ApiServer struct {
	httpSvr *http.Server
	JobMgr  JobOperator
}

var (
	G_ApiServer *ApiServer
)

func jobSaveHandler(resp http.ResponseWriter, req *http.Request) {
	bodyReader := req.Body
	job, err := common.UnMarshalJobFromReader(bodyReader)
	if err != nil {
		bytes, _ := buildResponse(-1, err.Error(), nil)
		resp.Write(bytes)
	}
	err = G_ApiServer.JobMgr.SaveJob(job)
	if err != nil {
		bytes, _ := buildResponse(-1, err.Error(), nil)
		resp.Write(bytes)
	}
	bytes, err := buildResponse(0, "suc", nil)
	resp.Write(bytes)
}

func JobDeleteHandler(resp http.ResponseWriter, req *http.Request) {
	bodyReader := req.Body
	job, err := common.UnMarshalJobFromReader(bodyReader)
	if err != nil {
		bytes, _ := buildResponse(-1, err.Error(), nil)
		resp.Write(bytes)
	}
	err = G_ApiServer.JobMgr.DeleteJob(job.Name)
	if err != nil {
		bytes, _ := buildResponse(-1, err.Error(), nil)
		resp.Write(bytes)
	}
	bytes, err := buildResponse(0, "delete suc", nil)
	resp.Write(bytes)
}

func JobListHandler(resp http.ResponseWriter, req *http.Request) {
	jobs, err := G_ApiServer.JobMgr.ListJobs()
	if err != nil {
		bytes, _ := buildResponse(-1, err.Error(), nil)
		resp.Write(bytes)
	}
	bytes, err := buildResponse(0, "list suc", jobs)
	resp.Write(bytes)
}

func JobkillHandler(resp http.ResponseWriter, req *http.Request) {
	bodyReader := req.Body
	defer req.Body.Close()
	job, err := common.UnMarshalJobFromReader(bodyReader)
	if err != nil {
		bytes, _ := buildResponse(-1, err.Error(), nil)
		resp.Write(bytes)
	}
	err = G_ApiServer.JobMgr.KillJob(job.Name)
	if err != nil {
		bytes, _ := buildResponse(-1, err.Error(), nil)
		resp.Write(bytes)
	}
	bytes, err := buildResponse(0, "kill suc", nil)
	resp.Write(bytes)
}

func testHandler(resp http.ResponseWriter, req *http.Request) {
	content := []byte("test ok")
	resp.Write(content)
}

func InitApiServer(j JobOperator) *ApiServer {
	//apiSvrPort := config.G_Config.ApiServerPort
	apiSvrRdTimeout := config.G_Config.ApiServerReadTimeOut
	apiSvrWtTimeout := config.G_Config.ApiServerWriteTimeOut

	if G_ApiServer != nil {
		return G_ApiServer
	}
	//route
	mux := http.NewServeMux()
	mux.HandleFunc("/job/test", testHandler)
	mux.HandleFunc("/job/save", jobSaveHandler)
	mux.HandleFunc("/job/list", JobListHandler)
	mux.HandleFunc("/job/delete", JobDeleteHandler)
	mux.HandleFunc("/job/kill", JobkillHandler)

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(config.G_Config.WebRoot))))

	httpSvr := &http.Server{
		ReadTimeout:  time.Duration(apiSvrRdTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(apiSvrWtTimeout) * time.Millisecond,
		Handler:      mux,
	}

	G_ApiServer = &ApiServer{httpSvr: httpSvr, JobMgr: j}
	return G_ApiServer
}

func (s *ApiServer) Run() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(config.G_Config.ApiServerPort))
	if err != nil {
		panic(err)
	}
	go s.httpSvr.Serve(listener)
}
