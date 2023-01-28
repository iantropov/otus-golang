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

	terminate := make(chan struct{})
	outs := make([]Bi, 0)

	isDone := false
	isClosed := false

	wg := sync.WaitGroup{}
	wg.Add(1)

	stagesWithDone := make([]Stage, 0, len(stages))
	for i := 0; i < len(stages); i++ {
		stagesWithDone = append(stagesWithDone, stageWithDone(done, stages[i]))
	}

	go func() {
		for {
			if !isClosed {
				select {
				case vIn := <-in:
					if vIn == nil {
						isClosed = true
						wg.Done()
						break
					}

					nextOut := make(Bi, 1)
					outs = append(outs, nextOut)

					wg.Add(1)
					go startPipelineWithValue(vIn, stagesWithDone, nextOut, &wg)
				case <-done:
					isDone = true
					return
				}
			} else {
				select {
				case <-terminate:
					return
				case <-done:
					isDone = true
					return
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		if !isDone {
			for _, nextOut := range outs {
				out <- <-nextOut
				close(nextOut)
			}
		}
		close(terminate)
		close(out)
	}()

	return out
}

func startPipelineWithValue(vIn interface{}, stages []Stage, out Bi, wg *sync.WaitGroup) {
	defer wg.Done()
	nextIn := make(Bi, 1)
	nextIn <- vIn
	vOut := passStages(nextIn, stages)
	if vOut == nil {
	} else {
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
