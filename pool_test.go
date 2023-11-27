package worker

import (
	"sync"
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

	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			val, err := pool.RunJob(j)
			assert.Equal(t, 0, val)
			assert.Nil(t, err)
		}()
	}
	wg.Wait()

	assert.Equal(t, 100, j.value.Lock().Get())
}
