package worker

type worker struct {
	contextChan <-chan *context
	stopChan    chan struct{}
}

func newWorker(contextChan <-chan *context, stopChan chan struct{}) *worker {
	return &worker{
		contextChan: contextChan,
		stopChan:    stopChan,
	}
}

func (w *worker) run() {
	for {
		select {
		case ctx := <-w.contextChan:
			w.runJob(ctx)
		case <-w.stopChan:
			return
		}
	}
}

func (w *worker) runJob(ctx *context) {
	val, err := ctx.job.Run()
	ctx.futureValue.Resolve(val, err)
}
