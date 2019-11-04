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
	contextChan  chan *context
	stopChan     chan struct{}
}

func NewPool(size PoolSize) Pool {
	contextChan := make(chan *context)
	stopChan := make(chan struct{})

	pool := pool{
		workersCount: size.getWorkersCount(),
		contextChan:  contextChan,
		stopChan:     stopChan,
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
	for i := 0; i < p.workersCount; i++ {
		go p.runWorker()
	}
}

func (p *pool) Stop() {
	for i := 0; i < p.workersCount; i++ {
		p.stopChan <- struct{}{}
	}
}

func (p *pool) runWorker() {
	for {
		select {
		case ctx := <-p.contextChan:
			p.runJob(ctx)
		case <-p.stopChan:
			return
		}
	}
}

func (p *pool) runJob(ctx *context) {
	val, err := ctx.job.Run()
	ctx.futureValue.Resolve(val, err)
}
