package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		return nil
	}
	nextIn := in
	for i := 0; i < len(stages); i++ {
		stageIn := make(Bi)
		go func(in In) {
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
		}(nextIn)
		nextIn = stages[i](stageIn)
	}
	return nextIn
}
