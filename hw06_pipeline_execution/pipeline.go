package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	cur := in

	for _, stage := range stages {
		cur = stage(orDone(done, cur))
	}

	return cur
}

func orDone(done In, in In) Out {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-done:
					for range in {
					}
					return
				}

			case <-done:
				for range in {
				}
				return
			}
		}
	}()

	return out
}
