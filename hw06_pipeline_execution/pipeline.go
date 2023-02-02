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
	stagesWithDone := make([]Stage, 0, len(stages))
	for i := 0; i < len(stages); i++ {
		stagesWithDone = append(stagesWithDone, stageWithDone(done, stages[i]))
	}
	wg := sync.WaitGroup{}
	for vIn := range in {
		wg.Add(1)
		go pipeline(vIn, done, stagesWithDone, out, &wg)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func pipeline(vIn interface{}, done In, stages []Stage, out Bi, wg *sync.WaitGroup) {
	defer wg.Done()
	in := make(Bi, 1)
	defer close(in)
	in <- vIn
	stageOut := stages[0](in)
	for i := 1; i < len(stages); i++ {
		stageOut = stages[i](stageOut)
	}
	select {
	case vOut := <-stageOut:
		if vOut != nil {
			out <- vOut
		}
	case <-done:
	}
}

func stageWithDone(done In, stage Stage) Stage {
	return func(in In) (out Out) {
		stageIn := make(Bi)
		go func() {
			defer close(stageIn)
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
