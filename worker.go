package worker

type worker interface {
	start()
	stop()
}

type workerImpl struct {
	contextChan <-chan *context
	stopChan    chan struct{}
}

func newWorker(contextChan <-chan *context) worker {
	return &workerImpl{
		contextChan: contextChan,
		stopChan:    make(chan struct{}),
	}
}

func (w *workerImpl) start() {
	go w.run()
}

func (w *workerImpl) run() {
	for {
		select {
		case ctx := <-w.contextChan:
			w.runJob(ctx)
		case <-w.stopChan:
			return
		}
	}
}

func (w *workerImpl) runJob(ctx *context) {
	err := ctx.job.Run()
	ctx.futureError.Resolve(err)
}

func (w *workerImpl) stop() {
	w.stopChan <- struct{}{}
}
