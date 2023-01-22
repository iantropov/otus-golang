package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("fewer tasks than the number of workers", func(t *testing.T) {
		tasksCount := 5
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 10
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, runTasksCount, int32(workersCount), "extra tasks were started")
	})
}

func TestRunForEdgeCases(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("empty tasks", func(t *testing.T) {
		tasksCount := 0
		tasks := make([]Task, 0, tasksCount)

		workersCount := 10
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
	})

	t.Run("empty workersCount", func(t *testing.T) {
		tasksCount := 5
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 0
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)

		require.Equal(t, int32(workersCount), runTasksCount, "no tasks were started")
	})

	t.Run("zero maxErrorsCount", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, int32(tasksCount), runTasksCount, "excessive tasks were started")
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	})

	t.Run("negative maxErrorsCount", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := -1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, int32(tasksCount), runTasksCount, "excessive tasks were started")
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	})
}

func TestRunConcurrency(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("test concurrency without Sleep", func(t *testing.T) {
		tasksCount := 10000
		tasks := make([]Task, 0, tasksCount)

		idxCh := make(chan int, tasksCount)
		for i := 0; i < tasksCount; i++ {
			idx := i
			tasks = append(tasks, func() error {
				idxCh <- idx
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		outOfOrder := false
		for i := 0; i < tasksCount; i++ {
			idx := <-idxCh
			if idx != i {
				outOfOrder = true
				break
			}
		}

		require.True(t, outOfOrder, "concurrent tasks should run out of linear order", err)
	})
}
