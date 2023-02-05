package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil || len(stages) == 0 {
		return nil
	}
	for _, stage := range stages {
		if stage == nil {
			return nil
		}
	}
	nextIn := in
	for _, stage := range stages {
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
		nextIn = stage(stageIn)
	}
	return nextIn
}
