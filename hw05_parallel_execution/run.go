package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	okChannel := make(chan bool, n)
	taskChannel := make(chan Task, n)
	wg := sync.WaitGroup{}

	errorsCounter := 0
	finishedCounter := 0
	tasksPointer := 0

	if n == 0 || len(tasks) == 0 {
		return nil
	}
	if m > n {
		m = n
	}

	for ; tasksPointer < n && tasksPointer < len(tasks); tasksPointer++ {
		taskChannel <- tasks[tasksPointer]
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(taskChannel, okChannel, &wg)
	}

	var err error
	for ok := range okChannel {
		if !ok && m > 0 {
			errorsCounter++
		}
		finishedCounter++

		if m > 0 && errorsCounter >= m {
			err = ErrErrorsLimitExceeded
			break
		} else if finishedCounter == len(tasks) {
			break
		}

		if tasksPointer < len(tasks) {
			taskChannel <- tasks[tasksPointer]
			tasksPointer++
		}
	}

	close(taskChannel)
	wg.Wait()
	close(okChannel)
	return err
}

func worker(taskChannel <-chan Task, okChannel chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskChannel {
		if task == nil {
			return
		}
		err := task()
		okChannel <- err == nil
		if err != nil {
			return
		}
	}
}
