// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/TeoPlow/online-music-service/src/musical/internal/domain (interfaces: AlbumRepo)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/TeoPlow/online-music-service/src/musical/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockAlbumRepo is a mock of AlbumRepo interface.
type MockAlbumRepo struct {
	ctrl     *gomock.Controller
	recorder *MockAlbumRepoMockRecorder
}

// MockAlbumRepoMockRecorder is the mock recorder for MockAlbumRepo.
type MockAlbumRepoMockRecorder struct {
	mock *MockAlbumRepo
}

// NewMockAlbumRepo creates a new mock instance.
func NewMockAlbumRepo(ctrl *gomock.Controller) *MockAlbumRepo {
	mock := &MockAlbumRepo{ctrl: ctrl}
	mock.recorder = &MockAlbumRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAlbumRepo) EXPECT() *MockAlbumRepoMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockAlbumRepo) Add(arg0 context.Context, arg1 models.Album) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockAlbumRepoMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockAlbumRepo)(nil).Add), arg0, arg1)
}

// Delete mocks base method.
func (m *MockAlbumRepo) Delete(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAlbumRepoMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAlbumRepo)(nil).Delete), arg0, arg1)
}

// FindCopy mocks base method.
func (m *MockAlbumRepo) FindCopy(arg0 context.Context, arg1 models.Album) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCopy", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCopy indicates an expected call of FindCopy.
func (mr *MockAlbumRepoMockRecorder) FindCopy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCopy", reflect.TypeOf((*MockAlbumRepo)(nil).FindCopy), arg0, arg1)
}

// GetByFilter mocks base method.
func (m *MockAlbumRepo) GetByFilter(arg0 context.Context, arg1, arg2 int, arg3 *string, arg4 *uuid.UUID) ([]models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByFilter", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByFilter indicates an expected call of GetByFilter.
func (mr *MockAlbumRepoMockRecorder) GetByFilter(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByFilter", reflect.TypeOf((*MockAlbumRepo)(nil).GetByFilter), arg0, arg1, arg2, arg3, arg4)
}

// GetByID mocks base method.
func (m *MockAlbumRepo) GetByID(arg0 context.Context, arg1 uuid.UUID) (models.Album, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(models.Album)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockAlbumRepoMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAlbumRepo)(nil).GetByID), arg0, arg1)
}

// Update mocks base method.
func (m *MockAlbumRepo) Update(arg0 context.Context, arg1 models.Album) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAlbumRepoMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAlbumRepo)(nil).Update), arg0, arg1)
}
