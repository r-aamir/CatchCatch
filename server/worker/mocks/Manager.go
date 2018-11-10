// Code generated by mockery v1.0.0
package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import worker "github.com/perenecabuto/CatchCatch/server/worker"

// Manager is an autogenerated mock type for the Manager type
type Manager struct {
	mock.Mock
}

// Add provides a mock function with given fields: w
func (_m *Manager) Add(w worker.Task) {
	_m.Called(w)
}

// BusyTasks provides a mock function with given fields:
func (_m *Manager) BusyTasks() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
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

// Flush provides a mock function with given fields:
func (_m *Manager) Flush() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields: w, params
func (_m *Manager) Run(w worker.Task, params worker.TaskParams) error {
	ret := _m.Called(w, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(worker.Task, worker.TaskParams) error); ok {
		r0 = rf(w, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunUnique provides a mock function with given fields: w, params
func (_m *Manager) RunUnique(w worker.Task, params worker.TaskParams) error {
	ret := _m.Called(w, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(worker.Task, worker.TaskParams) error); ok {
		r0 = rf(w, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RunningJobs provides a mock function with given fields:
func (_m *Manager) RunningJobs() ([]worker.Job, error) {
	ret := _m.Called()

	var r0 []worker.Job
	if rf, ok := ret.Get(0).(func() []worker.Job); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]worker.Job)
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

// Start provides a mock function with given fields: ctx
func (_m *Manager) Start(ctx context.Context) {
	_m.Called(ctx)
}

// Started provides a mock function with given fields:
func (_m *Manager) Started() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *Manager) Stop() {
	_m.Called()
}

// TasksID provides a mock function with given fields:
func (_m *Manager) TasksID() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}
