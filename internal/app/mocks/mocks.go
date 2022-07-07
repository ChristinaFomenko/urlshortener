// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	configs "github.com/ChristinaFomenko/shortener/configs"
	storage "github.com/ChristinaFomenko/shortener/internal/app/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddBatchURL mocks base method.
func (m *MockRepository) AddBatchURL(arg0 []storage.UserURL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBatchURL", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBatchURL indicates an expected call of AddBatchURL.
func (mr *MockRepositoryMockRecorder) AddBatchURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBatchURL", reflect.TypeOf((*MockRepository)(nil).AddBatchURL), arg0)
}

// AddURL mocks base method.
func (m *MockRepository) AddURL(userShortURL storage.UserURL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddURL", userShortURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddURL indicates an expected call of AddURL.
func (mr *MockRepositoryMockRecorder) AddURL(userShortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddURL", reflect.TypeOf((*MockRepository)(nil).AddURL), userShortURL)
}

// Destruct mocks base method.
func (m *MockRepository) Destruct(cfg configs.AppConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destruct", cfg)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destruct indicates an expected call of Destruct.
func (mr *MockRepositoryMockRecorder) Destruct(cfg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destruct", reflect.TypeOf((*MockRepository)(nil).Destruct), cfg)
}

// GetList mocks base method.
func (m *MockRepository) GetList(arg0 string) []storage.UserURL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", arg0)
	ret0, _ := ret[0].([]storage.UserURL)
	return ret0
}

// GetList indicates an expected call of GetList.
func (mr *MockRepositoryMockRecorder) GetList(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockRepository)(nil).GetList), arg0)
}

// GetShortByOriginal mocks base method.
func (m *MockRepository) GetShortByOriginal(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortByOriginal", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortByOriginal indicates an expected call of GetShortByOriginal.
func (mr *MockRepositoryMockRecorder) GetShortByOriginal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortByOriginal", reflect.TypeOf((*MockRepository)(nil).GetShortByOriginal), arg0)
}

// GetURL mocks base method.
func (m *MockRepository) GetURL(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockRepositoryMockRecorder) GetURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockRepository)(nil).GetURL), arg0)
}

// Ping mocks base method.
func (m *MockRepository) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockRepositoryMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockRepository)(nil).Ping))
}
