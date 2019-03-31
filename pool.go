package worker

type Pool interface {
	StartJob(Job) FutureError

	Start()
	Stop()
}

type pool struct {
	workers     []worker
	contextChan chan *context
}

func NewPool(size PoolSize) Pool {
	contextChan := make(chan *context)

	pool := pool{
		workers:     make([]worker, size.getWorkersCount()),
		contextChan: contextChan,
	}

	for index := range pool.workers {
		pool.workers[index] = newWorker(contextChan)
	}

	return &pool
}

func (p *pool) StartJob(job Job) FutureError {
	futureError := newFutureError()

	p.contextChan <- &context{
		job:         job,
		futureError: futureError,
	}

	return futureError
}

func (p *pool) Start() {
	for _, worker := range p.workers {
		worker.start()
	}
}

func (p *pool) Stop() {
	for _, worker := range p.workers {
		worker.stop()
	}
}
