package worker

import (
	"context"
	"strings"
	"time"

	"github.com/tidwall/sjson"
)

// Manager for workers and its tasks
type Manager interface {
	Start(ctx context.Context)
	Stop()
	Started() bool
	Add(w Task)

	Run(w Task, params TaskParams) error
	RunUnique(w Task, params TaskParams) error

	BusyTasks() ([]string, error)
	RunningJobs() ([]Job, error)
	Flush() error
}

// TaskParams is map of task parameters
type TaskParams map[string]interface{}

// Task runs tasks
type Task interface {
	ID() string
	Run(ctx context.Context, params TaskParams) error
}

// Job represents a worker job
type Job struct {
	ID         string
	TaskID     string
	Unique     bool
	Params     TaskParams
	Host       string
	LastUpdate time.Time
}

// IsUpdatedToInterval check if the job was update in a given interval
func (j *Job) IsUpdatedToInterval(interval time.Duration) bool {
	deadline := time.Now().Add(-interval)
	return j.LastUpdate.After(deadline)
}

// LockName return a unique lock name for this task
func (j *Job) LockName() string {
	params, _ := sjson.Set("", "params", j.Params)
	return strings.Join([]string{tasksQueue, j.TaskID, params, "lock"}, ":")
}
