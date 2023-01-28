package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	if in == nil {
		close(out)
		return out
	}

	isDone := make(chan struct{}, 1)
	outs := make([]Bi, 0)

	wgPipeline := sync.WaitGroup{}
	wgMaster := sync.WaitGroup{}
	wgWorkers := sync.WaitGroup{}
	wgMaster.Add(1)
	wgPipeline.Add(1)

	stagesWithDone := make([]Stage, 0, len(stages))
	for i := 0; i < len(stages); i++ {
		stagesWithDone = append(stagesWithDone, stageWithDone(done, stages[i]))
	}

	go func() {
		defer wgMaster.Done()
		for vIn := range in {
			nextOut := make(Bi, 1)
			outs = append(outs, nextOut)

			wgWorkers.Add(1)
			go pipelineWorker(vIn, stagesWithDone, nextOut, &wgWorkers)
		}
	}()

	go func() {
		wgMaster.Wait()
		wgWorkers.Wait()
		select {
		case <-done:
			isDone <- struct{}{}
		default:
		}
		wgPipeline.Done()
	}()

	go func() {
		defer func() {
			close(isDone)
			close(out)
		}()

		wgPipeline.Wait()

		select {
		case <-isDone:
			return
		default:
		}

		for _, nextOut := range outs {
			out <- <-nextOut
			close(nextOut)
		}
	}()

	return out
}

func pipelineWorker(vIn interface{}, stages []Stage, out Bi, wg *sync.WaitGroup) {
	defer wg.Done()
	nextIn := make(Bi, 1)
	nextIn <- vIn
	vOut := passStages(nextIn, stages)
	if vOut != nil {
		out <- vOut
	}
}

func passStages(in In, stages []Stage) interface{} {
	nextIn := in
	var lastOut Out
	for _, stage := range stages {
		lastOut = stage(nextIn)
		nextIn = lastOut
	}
	return <-lastOut
}

func stageWithDone(done In, stage Stage) Stage {
	return func(in In) (out Out) {
		stageIn := make(Bi)

		go func() {
			defer func() {
				close(stageIn)
			}()

			for {
				select {
				case vIn := <-in:
					if vIn == nil {
						return
					}
					stageIn <- vIn
				case <-done:
					return
				}
			}
		}()

		return stage(stageIn)
	}
}
