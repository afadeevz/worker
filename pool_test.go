package worker

import (
	"testing"

	"github.com/AlexanderFadeev/lock"
	"github.com/stretchr/testify/assert"
)

type job struct {
	value lock.Value
}

func (j *job) Run() (int, error) {
	value := j.value.Lock()
	defer value.Unlock()

	value.Set(value.Get().(int) + 1)
	return 0, nil
}

func TestPool(t *testing.T) {
	pool := NewPool[int](16)
	defer pool.Stop()

	j := &job{
		value: lock.NewValue(0),
	}

	for i := 0; i < 100; i++ {
		go func() {
			val, err := pool.RunJob(j)
			assert.Equal(t, 0, val)
			assert.Nil(t, err)
		}()
	}

	assert.Equal(t, 100, j.value.Lock().Get())
}
