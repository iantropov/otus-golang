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

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	wg := sync.WaitGroup{}
	wg.Add(1)

	counter := 0
	outs := make([]Out, 0)
	go func() {
		defer wg.Done()
		for {
			select {
			case vIn, ok := <-in:
				fmt.Println("Received VALUE #", counter, ":", vIn)
				if !ok {
					fmt.Println("Received CLOSED")
					return
				}

				wg.Add(1)
				nextOut := make(Bi, 1)
				outs = append(outs, nextOut)
				go func(counter int, out Bi) {
					defer wg.Done()
					nextIn := make(Bi, 1)
					nextIn <- vIn
					fmt.Println("Send VALUE #", counter, "into the pipe (", vIn, ")")
					vOut := passStages(nextIn, stages)
					fmt.Println("Recevied RESULT #", counter, "from the pipe (", vOut, ")")
					out <- vOut
				}(counter, nextOut)
				counter++
			case <-done:
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		fmt.Println("READY TO OUTPUT RESULTS")
		for _, nextOut := range outs {
			out <- <-nextOut
		}
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
