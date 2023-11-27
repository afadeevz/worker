package worker

type Job[T any] interface {
	Run() (T, error)
}

type Func[T any] func() (T, error)

func (f Func[T]) Run() (T, error) {
	return f()
}
