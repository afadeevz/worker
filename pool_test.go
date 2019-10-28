package worker

import (
	"github.com/AlexanderFadeev/future"
	"github.com/AlexanderFadeev/lock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type job struct {
	value lock.Value
}

func (j *job) Run() (interface{}, error) {
	value := j.value.Lock()
	defer value.Unlock()

	value.Set(value.Get().(int) + 1)
	return nil, nil
}

func TestPool(t *testing.T) {
	pool := NewPool(16)
	pool.Start()
	defer pool.Stop()

	j := &job{
		value: lock.NewValue(0),
	}

	var futures []future.Future
	for i := 0; i < 100; i++ {
		f := pool.StartJob(j)
		futures = append(futures, f)
	}

	for _, f := range futures {
		val, err := f.Wait()
		assert.Nil(t, val)
		assert.Nil(t, err)
	}
	assert.Equal(t, j.value.Lock().Get(), 100)
}
