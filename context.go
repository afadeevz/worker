package worker

type context struct {
	job         Job
	futureError *futureError
}
