package redis

type InputData interface{}
type Result interface{}

type Work struct {
	Workable func(InputData) Result
	Input    InputData
	Result   Result
}

type WorkerPool struct {
	Input  chan *Work
	Output chan *Work
}

func NewWorkerPool(workers int) WorkerPool {
	pool := WorkerPool{
		Input:  make(chan *Work),
		Output: make(chan *Work),
	}
	for i := 0; i < workers; i++ {
		go pool.work()
	}
	return pool
}

func (pool WorkerPool) work() {
	for {
		w := <-pool.Input
		w.Result = w.Workable(w.Input)
		pool.Output <- w
	}
}