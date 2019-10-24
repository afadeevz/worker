package worker

import "github.com/AlexanderFadeev/myerrors"

type Pool interface {
	StartJob(Job) FutureError
	RunJobs(JobGenerator) error

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

func (p *pool) RunJobs(jg JobGenerator) error {
	var futureErrs []FutureError

	for job := jg.GetJob(); job != nil; job = jg.GetJob() {
		futureErr := p.StartJob(job)
		futureErrs = append(futureErrs, futureErr)
	}

	var err error
	for _, fErr := range futureErrs {
		err = myerrors.Merge(err, fErr.Get())
	}

	return err
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
