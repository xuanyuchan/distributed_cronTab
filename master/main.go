package main

import (
	"distributed_cronTab/master/Api"
	"distributed_cronTab/master/config"
	"distributed_cronTab/master/manager"
	"fmt"
)

func initConfig() {
	err := config.ParseConfig("./master/master.json")
	if err != nil {
		panic(err)
	}
}

func initApiServer(operator Api.JobOperator) {
	Api.InitApiServer(operator).Run()
}

func initEtcdClient() *manager.JobManager {
	jobMgr, err := manager.InitJobManager()
	if err != nil {
		panic(err)
	}
	return jobMgr
}

func main() {
	fmt.Println("master start")
	ch := make(chan struct{})

	initConfig()

	//build etcd client
	jm := initEtcdClient()

	//launch api server
	initApiServer(jm)

	<-ch
}
