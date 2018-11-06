// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import time "time"
import worker "github.com/perenecabuto/CatchCatch/server/worker"

// TaskManagerQueue is an autogenerated mock type for the TaskManagerQueue type
type TaskManagerQueue struct {
	mock.Mock
}

// EnqueuePending provides a mock function with given fields: _a0
func (_m *TaskManagerQueue) EnqueuePending(_a0 *worker.Job) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*worker.Job) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EnqueueToProcess provides a mock function with given fields: _a0
func (_m *TaskManagerQueue) EnqueueToProcess(_a0 *worker.Job) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*worker.Job) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetJobByID provides a mock function with given fields: _a0
func (_m *TaskManagerQueue) GetJobByID(_a0 string) (*worker.Job, error) {
	ret := _m.Called(_a0)

	var r0 *worker.Job
	if rf, ok := ret.Get(0).(func(string) *worker.Job); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*worker.Job)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HeartbeatJob provides a mock function with given fields: _a0, _a1
func (_m *TaskManagerQueue) HeartbeatJob(_a0 *worker.Job, _a1 time.Duration) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*worker.Job, time.Duration) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// JobsOnProcessQueue provides a mock function with given fields:
func (_m *TaskManagerQueue) JobsOnProcessQueue() ([]*worker.Job, error) {
	ret := _m.Called()

	var r0 []*worker.Job
	if rf, ok := ret.Get(0).(func() []*worker.Job); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*worker.Job)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PollPending provides a mock function with given fields:
func (_m *TaskManagerQueue) PollPending() (*worker.Job, error) {
	ret := _m.Called()

	var r0 *worker.Job
	if rf, ok := ret.Get(0).(func() *worker.Job); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*worker.Job)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PollProcess provides a mock function with given fields:
func (_m *TaskManagerQueue) PollProcess() (*worker.Job, error) {
	ret := _m.Called()

	var r0 *worker.Job
	if rf, ok := ret.Get(0).(func() *worker.Job); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*worker.Job)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveFromProcessingQueue provides a mock function with given fields: _a0
func (_m *TaskManagerQueue) RemoveFromProcessingQueue(_a0 *worker.Job) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*worker.Job) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveJob provides a mock function with given fields: _a0
func (_m *TaskManagerQueue) RemoveJob(_a0 *worker.Job) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*worker.Job) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetJob provides a mock function with given fields: _a0
func (_m *TaskManagerQueue) SetJob(_a0 *worker.Job) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*worker.Job) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetJobLock provides a mock function with given fields: _a0, _a1
func (_m *TaskManagerQueue) SetJobLock(_a0 *worker.Job, _a1 time.Duration) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*worker.Job, time.Duration) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*worker.Job, time.Duration) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateJobLock provides a mock function with given fields: _a0, _a1
func (_m *TaskManagerQueue) UpdateJobLock(_a0 *worker.Job, _a1 time.Duration) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*worker.Job, time.Duration) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*worker.Job, time.Duration) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
