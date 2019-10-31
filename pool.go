package worker

import (
	"github.com/AlexanderFadeev/future"
)

type Pool interface {
	StartJob(Job) future.Future
	RunJob(Job) (interface{}, error)
	RunJobs(JobGenerator) ([]interface{}, error)

	Start()
	Stop()
}

type pool struct {
	workers     []*worker
	contextChan chan<- *context
	stopChan    chan<- struct{}
}

func NewPool(size PoolSize) Pool {
	contextChan := make(chan *context)
	stopChan := make(chan struct{})

	pool := pool{
		workers:     make([]*worker, size.getWorkersCount()),
		contextChan: contextChan,
		stopChan:    stopChan,
	}

	for index := range pool.workers {
		pool.workers[index] = newWorker(contextChan, stopChan)
	}

	return &pool
}

func (p *pool) StartJob(job Job) future.Future {
	f := future.NewValue()

	p.contextChan <- &context{
		job:         job,
		futureValue: f,
	}

	return f
}

func (p *pool) RunJob(job Job) (interface{}, error) {
	return p.StartJob(job).Wait()
}

func (p *pool) RunJobs(jg JobGenerator) (results []interface{}, err error) {
	var futures []future.Future

	for job := jg.GetJob(); job != nil; job = jg.GetJob() {
		f := p.StartJob(job)
		futures = append(futures, f)
	}

	for _, f := range futures {
		val, err := f.Wait()
		if err != nil {
			return nil, err
		}

		results = append(results, val)
	}

	return
}

func (p *pool) Start() {
	for _, worker := range p.workers {
		go worker.run()
	}
}

func (p *pool) Stop() {
	for range p.workers {
		p.stopChan <- struct{}{}
	}
}
