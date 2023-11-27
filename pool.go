package worker

import "github.com/AlexanderFadeev/future/v2"

type Pool[T any] interface {
	RunJob(Job[T]) (T, error)
	Stop()
}

type pool[T any] struct {
	funcChan chan func()
}

func NewPool[T any](size PoolSize) Pool[T] {
	p := &pool[T]{
		funcChan: make(chan func()),
	}

	for i := 0; i < size.getWorkersCount(); i++ {
		go p.runWorker()
	}

	return p
}

func (p *pool[T]) Stop() {
	close(p.funcChan)
}

func (p *pool[T]) runWorker() {
	for fn := range p.funcChan {
		fn()
	}
}

func (p *pool[T]) RunJob(job Job[T]) (result T, err error) {
	w := future.NewWaiter()

	p.funcChan <- func() {
		defer w.Done()
		result, err = job.Run()
	}

	w.Wait()
	return
}
