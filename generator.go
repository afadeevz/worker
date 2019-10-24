package worker

type JobGenerator interface {
	GetJob() Job
}
