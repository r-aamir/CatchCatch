package worker

import "context"

// Manager for workers and its tasks
type Manager interface {
	Start(ctx context.Context)
	Stop()
	Add(w Worker)
	Run(w Worker, params map[string]string) error
}

// Worker runs tasks
type Worker interface {
	ID() string
	Job(params map[string]string) error
}
