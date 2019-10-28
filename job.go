package worker

type Job interface {
	Run() (interface{}, error)
}

type Func func() (interface{}, error)

func (f Func) Run() (interface{}, error) {
	return f()
}
