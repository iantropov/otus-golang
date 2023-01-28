package hw06pipelineexecution

import (
	"fmt"
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)
type StageWithDone func(in In, done In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	wg := sync.WaitGroup{}
	wg.Add(1)
	terminate := make(chan struct{})
	counter := 0
	outs := make([]Out, 0)
	isDone := false
	isClosed := false
	// myClose := make(Out)

	stagesWithDone := make([]Stage, 0, len(stages))
	for i := 0; i < len(stages); i++ {
		stagesWithDone = append(stagesWithDone, stageWithDone(done, stages[i]))
	}

	go func() {
		for {
			if !isClosed {
				select {
				case vIn := <-in:
					fmt.Println("Received VALUE #", counter, ":", vIn)
					if vIn == nil {
						fmt.Println("Received CLOSED")
						isClosed = true
						wg.Done()
						break
					}

					wg.Add(1)
					nextOut := make(Bi, 1)
					outs = append(outs, nextOut)
					go func(counter int, out Bi) {
						defer wg.Done()
						nextIn := make(Bi, 1)
						nextIn <- vIn
						fmt.Println("Send VALUE #", counter, "into the pipe (", vIn, ")")
						vOut := passStages(nextIn, stagesWithDone)
						fmt.Println("Recevied RESULT #", counter, "from the pipe (", vOut, ")")
						if vOut == nil {
							fmt.Println("Received CLOSED RESULT")
						} else {
							out <- vOut
						}
					}(counter, nextOut)
					counter++
				case <-done:
					fmt.Println("Received DONE")
					isDone = true
					return
				}
			} else {
				select {
				case <-terminate:
					return
				case <-done:
					fmt.Println("Received DONE2")
					isDone = true
					return
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		fmt.Println("READY TO OUTPUT RESULTS")
		if !isDone {
			for _, nextOut := range outs {
				out <- <-nextOut
			}
		}
		close(terminate)
		close(out)
	}()

	return out
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
				fmt.Println("CLOSED INTER CHANNEL")
				close(stageIn)
			}()

			fmt.Println("STARTED INTER CHANNEL")

			for {
				select {
				case vIn := <-in:
					fmt.Println("INTER VALUE", vIn)
					if vIn == nil {
						fmt.Println("Received INTER CLOSED")
						return
					}
					stageIn <- vIn
				case <-done:
					fmt.Println("INTER DONE")
					return
				}
			}
		}()

		return stage(stageIn)
	}
}
