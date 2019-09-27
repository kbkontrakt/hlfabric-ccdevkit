// Code generated by MockGen. DO NOT EDIT.
// Source: time_svc.go

// Package chrono is a generated GoMock package.
package chrono

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTimeService is a mock of TimeService interface
type MockTimeService struct {
	ctrl     *gomock.Controller
	recorder *MockTimeServiceMockRecorder
}

// MockTimeServiceMockRecorder is the mock recorder for MockTimeService
type MockTimeServiceMockRecorder struct {
	mock *MockTimeService
}

// NewMockTimeService creates a new mock instance
func NewMockTimeService(ctrl *gomock.Controller) *MockTimeService {
	mock := &MockTimeService{ctrl: ctrl}
	mock.recorder = &MockTimeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTimeService) EXPECT() *MockTimeServiceMockRecorder {
	return m.recorder
}

// NowDateTime mocks base method
func (m *MockTimeService) NowDateTime() (DateTime, error) {
	ret := m.ctrl.Call(m, "NowDateTime")
	ret0, _ := ret[0].(DateTime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NowDateTime indicates an expected call of NowDateTime
func (mr *MockTimeServiceMockRecorder) NowDateTime() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NowDateTime", reflect.TypeOf((*MockTimeService)(nil).NowDateTime))
}
