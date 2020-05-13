package workerpool

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestInitializeActivePool(t *testing.T) {
	tp := InitializeActivePool(3)
	expectationsMatch := tp.poolSize == 3 && cap(tp.tRlt) == 3 && cap(tp.tRgt) == 3
	if !expectationsMatch {
		t.Error()
	}
}

func TestActivePool_AddWork(t *testing.T) {
	tp := InitializeActivePool(4)
	for i := 0; i < 5; i++ {
		wrk := func(tskRslt TaskResult) {
			tskRslt <- nil
		}
		tp.AddWork(wrk)
	}
	expectationsMatch := tp.poolSize == 4 && cap(tp.tRlt) == 4 && cap(tp.tRgt) == 4 && len(tp.taskList) == 5
	if !expectationsMatch {
		t.Error()
	}
}

func TestActivePool_DoWork(t *testing.T) {
	tp := InitializeActivePool(2)
	for i := 0; i < 10; i++ {
		wrk := func(tskRslt TaskResult) {
			sleepDuration := time.Duration(rand.Intn(10))
			time.Sleep(sleepDuration * time.Millisecond)
			fmt.Printf("Sleeping for ... %v\n", sleepDuration)
			tskRslt <- nil
		}
		tp.AddWork(wrk)
	}
	resultCount, successCount := tp.DoWork()
	if resultCount != 10 || successCount != 10 {
		t.Error()
	}
}

func TestActivePool_SetTimeout(t *testing.T) {
	tp := InitializeActivePool(2)
	tp.SetTimeout(time.Duration(1))
	for i := 0; i < 10; i++ {
		wrk := func(tskRslt TaskResult) {
			time.Sleep(1 * time.Minute)
			tskRslt <- nil
		}
		tp.AddWork(wrk)
	}
	resultCount, successCount := tp.DoWork()
	if resultCount != 0 || successCount != 0 {
		t.Error()
	}
}
