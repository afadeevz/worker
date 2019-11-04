package worker

import "testing"

type sqrJob int64

func (j sqrJob) Run() (interface{}, error) {
	return j * j, nil
}

type sqrJobGenerator struct {
	i     int64
	limit int64
}

func (g *sqrJobGenerator) GetJob() Job {
	if g.i > g.limit {
		return nil
	}

	g.i++
	return sqrJob(g.i - 1)
}

func BenchmarkPool(b *testing.B) {
	pool := NewPool(SizeAuto)
	pool.Start()
	defer pool.Stop()

	generator := &sqrJobGenerator{
		i:     0,
		limit: int64(b.N),
	}

	_, _ = pool.RunJobs(generator)
}
