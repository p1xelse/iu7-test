// Code generated by MockGen. DO NOT EDIT.
// Source: metrics.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetrics is a mock of Metrics interface.
type MockMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsMockRecorder
}

// MockMetricsMockRecorder is the mock recorder for MockMetrics.
type MockMetricsMockRecorder struct {
	mock *MockMetrics
}

// NewMockMetrics creates a new mock instance.
func NewMockMetrics(ctrl *gomock.Controller) *MockMetrics {
	mock := &MockMetrics{ctrl: ctrl}
	mock.recorder = &MockMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetrics) EXPECT() *MockMetricsMockRecorder {
	return m.recorder
}

// DecodeRecordTime mocks base method.
func (m *MockMetrics) DecodeRecordTime(decodingAlgorithm string, duration float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DecodeRecordTime", decodingAlgorithm, duration)
}

// DecodeRecordTime indicates an expected call of DecodeRecordTime.
func (mr *MockMetricsMockRecorder) DecodeRecordTime(decodingAlgorithm, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodeRecordTime", reflect.TypeOf((*MockMetrics)(nil).DecodeRecordTime), decodingAlgorithm, duration)
}

// EncodeRecordTime mocks base method.
func (m *MockMetrics) EncodeRecordTime(encodingAlgorithm string, duration float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EncodeRecordTime", encodingAlgorithm, duration)
}

// EncodeRecordTime indicates an expected call of EncodeRecordTime.
func (mr *MockMetricsMockRecorder) EncodeRecordTime(encodingAlgorithm, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncodeRecordTime", reflect.TypeOf((*MockMetrics)(nil).EncodeRecordTime), encodingAlgorithm, duration)
}