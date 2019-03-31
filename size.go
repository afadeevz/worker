package worker

import (
	"runtime"
)

type PoolSize int

const SizeAuto PoolSize = 0

func (s PoolSize) getWorkersCount() int {
	switch {
	case s <= 0:
		return runtime.NumCPU()
	default:
		return int(s)
	}
}
