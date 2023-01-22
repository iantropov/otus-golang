package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	okChannel := make(chan bool, len(tasks))
	taskChannel := make(chan Task, len(tasks))
	stopChannel := make(chan struct{})
	wg := sync.WaitGroup{}

	errorsCounter := 0
	finishedCounter := 0

	if n == 0 || len(tasks) == 0 {
		return nil
	}

	for _, task := range tasks {
		taskChannel <- task
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(taskChannel, okChannel, stopChannel, &wg)
	}

	var err error
	if m <= 0 {
		err = ErrErrorsLimitExceeded
	}
	for {
		ok := <-okChannel
		if !ok && m > 0 {
			errorsCounter++
		}
		finishedCounter++

		if m > 0 && errorsCounter >= m {
			err = ErrErrorsLimitExceeded
			for i := 0; i < n; i++ {
				stopChannel <- struct{}{}
			}
			break
		} else if finishedCounter == len(tasks) {
			break
		}
	}

	wg.Wait()

	return err
}

func worker(taskChannel <-chan Task, okChannel chan<- bool, stopChannel <-chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-stopChannel:
			wg.Done()
			return
		default:
		}
		select {
		case task := <-taskChannel:
			err := task()
			okChannel <- err == nil
		default:
			wg.Done()
			return
		}
	}
}
