package worker

import (
	"github.com/AlexanderFadeev/future"
)

type FutureError interface {
	Get() error
}

type futureErrorResolver interface {
	FutureError

	Resolve(err error)
}

type futureError struct {
	impl future.Resolver
}

func newFutureError() futureErrorResolver {
	return &futureError{
		impl: future.NewFuture(),
	}
}

func (f *futureError) Resolve(err error) {
	f.impl.Resolve(err)
}

func (f *futureError) Get() error {
	switch val := f.impl.Value(); val {
	case nil:
		return nil
	default:
		return val.(error)
	}
}
