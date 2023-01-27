package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// outs := make([]Out, len(stages))
	nextIn := in
	var lastOut Out
	for _, stage := range stages {
		lastOut = stage(nextIn)
		nextIn = lastOut
	}

	// out := fanIn(done, outs...)

	return lastOut
}

// func fanIn(done In, outs ...Out) Out {
// 	if len(outs) == 1 {
// 		return outs[0]
// 	}

// 	fanInOut := make(Bi)

// 	go func() {
// 		for {
// 			select {
// 			case o := <-outs[0]:
// 				fanInOut <- o
// 			case o := <-fanIn(done, outs[1:]...):
// 				fanInOut <- o
// 			case <-done:
// 				return
// 			}
// 		}
// 	}()

// 	return fanInOut
// }
