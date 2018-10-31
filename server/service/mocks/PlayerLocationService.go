// Code generated by mockery v1.0.0
package mocks

import context "context"
import mock "github.com/stretchr/testify/mock"
import model "github.com/perenecabuto/CatchCatch/server/model"
import service "github.com/perenecabuto/CatchCatch/server/service"

// PlayerLocationService is an autogenerated mock type for the PlayerLocationService type
type PlayerLocationService struct {
	mock.Mock
}

// All provides a mock function with given fields:
func (_m *PlayerLocationService) All() (model.PlayerList, error) {
	ret := _m.Called()

	var r0 model.PlayerList
	if rf, ok := ret.Get(0).(func() model.PlayerList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.PlayerList)
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

// Clear provides a mock function with given fields:
func (_m *PlayerLocationService) Clear() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Features provides a mock function with given fields:
func (_m *PlayerLocationService) Features() ([]*model.Feature, error) {
	ret := _m.Called()

	var r0 []*model.Feature
	if rf, ok := ret.Get(0).(func() []*model.Feature); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Feature)
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

// GeofenceByID provides a mock function with given fields: id
func (_m *PlayerLocationService) GeofenceByID(id string) (*model.Feature, error) {
	ret := _m.Called(id)

	var r0 *model.Feature
	if rf, ok := ret.Get(0).(func(string) *model.Feature); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Feature)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ObserveFeaturesEventsNearToAdmin provides a mock function with given fields: ctx, cb
func (_m *PlayerLocationService) ObserveFeaturesEventsNearToAdmin(ctx context.Context, cb service.AdminNearToFeatureCallback) error {
	ret := _m.Called(ctx, cb)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, service.AdminNearToFeatureCallback) error); ok {
		r0 = rf(ctx, cb)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ObservePlayerNearToCheckpoint provides a mock function with given fields: _a0, _a1
func (_m *PlayerLocationService) ObservePlayerNearToCheckpoint(_a0 context.Context, _a1 service.PlayerNearToFeatureCallback) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, service.PlayerNearToFeatureCallback) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ObservePlayersNearToGeofence provides a mock function with given fields: ctx, cb
func (_m *PlayerLocationService) ObservePlayersNearToGeofence(ctx context.Context, cb func(string, model.Player) error) error {
	ret := _m.Called(ctx, cb)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(string, model.Player) error) error); ok {
		r0 = rf(ctx, cb)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Remove provides a mock function with given fields: playerID
func (_m *PlayerLocationService) Remove(playerID string) error {
	ret := _m.Called(playerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(playerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveAdmin provides a mock function with given fields: id
func (_m *PlayerLocationService) RemoveAdmin(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Set provides a mock function with given fields: p
func (_m *PlayerLocationService) Set(p *model.Player) error {
	ret := _m.Called(p)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Player) error); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetAdmin provides a mock function with given fields: id, lat, lon
func (_m *PlayerLocationService) SetAdmin(id string, lat float64, lon float64) error {
	ret := _m.Called(id, lat, lon)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, float64, float64) error); ok {
		r0 = rf(id, lat, lon)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetCheckpoint provides a mock function with given fields: id, coordinates
func (_m *PlayerLocationService) SetCheckpoint(id string, coordinates string) error {
	ret := _m.Called(id, coordinates)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(id, coordinates)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetGeofence provides a mock function with given fields: id, coordinates
func (_m *PlayerLocationService) SetGeofence(id string, coordinates string) error {
	ret := _m.Called(id, coordinates)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(id, coordinates)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
