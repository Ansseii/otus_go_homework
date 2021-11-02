package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskCh := make(chan Task)
	var errorCount int32
	wg := sync.WaitGroup{}

	producer := func() {
		defer close(taskCh)
		for _, task := range tasks {
			if atomic.LoadInt32(&errorCount) >= int32(m) {
				break
			}
			taskCh <- task
		}
	}

	worker := func() {
		defer wg.Done()
		for task := range taskCh {
			if err := task(); err != nil {
				atomic.AddInt32(&errorCount, 1)
			}
		}
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker()
	}

	producer()
	wg.Wait()

	if errorCount >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
