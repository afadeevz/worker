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
	workersCount int
	funcChan     chan func()
}

func NewPool(size PoolSize) Pool {
	return &pool{
		workersCount: size.getWorkersCount(),
	}
}

func (p *pool) StartJob(job Job) future.Future {
	f := future.NewValue()

	p.funcChan <- func() {
		f.Resolve(job.Run())
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
	p.funcChan = make(chan func())
	for i := 0; i < p.workersCount; i++ {
		go p.runWorker()
	}
}

func (p *pool) Stop() {
	close(p.funcChan)
}

func (p *pool) runWorker() {
	for fn := range p.funcChan {
		fn()
	}
}
