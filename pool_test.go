package worker

import (
	"github.com/AlexanderFadeev/lock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type job struct {
	value lock.Value
}

func (j *job) Run() error {
	value := j.value.Lock()
	defer value.Unlock()

	value.Set(value.Get().(int) + 1)
	return nil
}

func TestPool(t *testing.T) {
	pool := NewPool(SizeAuto)
	pool.Start()
	defer pool.Stop()

	j := &job{
		value: lock.NewValue(0),
	}

	var errs []FutureError
	for i := 0; i < 100; i++ {
		err := pool.StartJob(j)
		errs = append(errs, err)
	}

	for _, err := range errs {
		assert.Nil(t, err.Get())
	}
	assert.Equal(t, j.value.Lock().Get(), 100)
}
