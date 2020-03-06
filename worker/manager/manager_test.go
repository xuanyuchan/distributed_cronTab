package manager

import (
	"distributed_cronTab/worker/config"
	"log"
	"testing"
)

func TestJobManager_WatchJobs(t *testing.T) {
	_ = config.ParseConfig("../worker.json")
	jM, err := InitJobManager()
	if err != nil {
		t.Errorf("init job manager error: %v\n", err)
		return
	}

	ch, err := jM.WatchJobs()
	if err != nil {
		t.Errorf("job watcher error: %v\n", err)
		return
	}
	for e := range ch {
		log.Printf("event, type:%d, job:%+v\n", e.EventType, e.Job)
	}

}
