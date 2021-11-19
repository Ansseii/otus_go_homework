package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func wrapper(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		out = stage(wrapper(out, done))
	}

	return out
}
