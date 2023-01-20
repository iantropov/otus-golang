package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	okChannel := make(chan bool)
	errorsCounter := 0
	runningTasks := 0
	tasksPointer := 0

	if m < 0 {
		m = 0
	}

	for i := 0; i < n && tasksPointer < len(tasks); i++ {
		go doTask(tasks[tasksPointer], okChannel)
		runningTasks++
		tasksPointer++
	}

	if runningTasks == 0 {
		return nil
	}

	var err error
	for isFinished := false; !isFinished; {
		select {
		case ok := <-okChannel:
			if !ok {
				errorsCounter++
			}

			if errorsCounter > m {
				err = ErrErrorsLimitExceeded
			}

			runningTasks--
			if runningTasks == 0 && (tasksPointer == len(tasks) || errorsCounter > m) {
				isFinished = true
				break
			}

			if runningTasks < n && errorsCounter <= m && tasksPointer < len(tasks) {
				go doTask(tasks[tasksPointer], okChannel)
				tasksPointer++
				runningTasks++
			}
		}
	}

	return err
}

func doTask(task Task, okChannel chan<- bool) {
	err := task()
	okChannel <- err == nil
}
