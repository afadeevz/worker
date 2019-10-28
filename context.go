package worker

import "github.com/AlexanderFadeev/future"

type context struct {
	job         Job
	futureValue future.Value
}
