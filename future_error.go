package worker

type FutureError interface {
	Resolve(err error)
	Get() error
}

type futureError struct {
	result chan error
}

func newFutureError() FutureError {
	return &futureError{
		result: make(chan error, 1),
	}
}

func (f *futureError) Resolve(err error) {
	f.result <- err
}

func (f *futureError) Get() error {
	result := <-f.result
	f.result <- result
	return result
}
