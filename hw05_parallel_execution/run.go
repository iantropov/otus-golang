package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	okChannel := make(chan bool, n)
	errorsCounter := 0
	runningTasks := 0
	tasksPointer := 0

	for i := 0; i < n && tasksPointer < len(tasks); i++ {
		go doTask(tasks[tasksPointer], okChannel)
		runningTasks++
		tasksPointer++
	}

	if runningTasks == 0 {
		return nil
	}

	var err error
	if m <= 0 {
		m = -1
		err = ErrErrorsLimitExceeded
	}

	for {
		ok := <-okChannel
		if !ok {
			errorsCounter++
		}

		if err == nil && errorsCounter > m {
			err = ErrErrorsLimitExceeded
		}

		runningTasks--
		if runningTasks == 0 && (tasksPointer == len(tasks) || errorsCounter > m) {
			break
		}

		if runningTasks < n && errorsCounter <= m && tasksPointer < len(tasks) {
			go doTask(tasks[tasksPointer], okChannel)
			tasksPointer++
			runningTasks++
		}
	}

	return err
}

func doTask(task Task, okChannel chan<- bool) {
	err := task()
	okChannel <- err == nil
}
