package util

import (
	"log"
	"sync"
	"testing"
)

func TestLock_Lock(t *testing.T) {
	sum := 0
	failNum := 0

	key := "/cron/test/testLock"
	lock, _ := InitLock(key)
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(wg *sync.WaitGroup) {
			err := lock.Lock()
			defer lock.UnLock()
			if err != nil {
				//log.Printf("lock fail\n")
				failNum += 1
			} else {
				log.Printf("lock suc\n")
				sum += 1
			}
			wg.Done()
		}(&wg)
	}
	wg.Wait()
	t.Logf("sum:%d, fail num:%d\n", sum, failNum)
}
