// Code generated by MockGen. DO NOT EDIT.
// Source: src/musical/internal/domain/like_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/TeoPlow/online-music-service/src/musical/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockLikeRepo is a mock of LikeRepo interface.
type MockLikeRepo struct {
	ctrl     *gomock.Controller
	recorder *MockLikeRepoMockRecorder
}

// MockLikeRepoMockRecorder is the mock recorder for MockLikeRepo.
type MockLikeRepoMockRecorder struct {
	mock *MockLikeRepo
}

// NewMockLikeRepo creates a new mock instance.
func NewMockLikeRepo(ctrl *gomock.Controller) *MockLikeRepo {
	mock := &MockLikeRepo{ctrl: ctrl}
	mock.recorder = &MockLikeRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLikeRepo) EXPECT() *MockLikeRepoMockRecorder {
	return m.recorder
}

// GetLikedArtists mocks base method.
func (m *MockLikeRepo) GetLikedArtists(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedArtists", ctx, userID, offset, limit)
	ret0, _ := ret[0].([]models.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedArtists indicates an expected call of GetLikedArtists.
func (mr *MockLikeRepoMockRecorder) GetLikedArtists(ctx, userID, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedArtists", reflect.TypeOf((*MockLikeRepo)(nil).GetLikedArtists), ctx, userID, offset, limit)
}

// GetLikedTracks mocks base method.
func (m *MockLikeRepo) GetLikedTracks(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLikedTracks", ctx, userID, offset, limit)
	ret0, _ := ret[0].([]models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLikedTracks indicates an expected call of GetLikedTracks.
func (mr *MockLikeRepoMockRecorder) GetLikedTracks(ctx, userID, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLikedTracks", reflect.TypeOf((*MockLikeRepo)(nil).GetLikedTracks), ctx, userID, offset, limit)
}

// LikeArtist mocks base method.
func (m *MockLikeRepo) LikeArtist(ctx context.Context, userID, artistID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikeArtist", ctx, userID, artistID)
	ret0, _ := ret[0].(error)
	return ret0
}

// LikeArtist indicates an expected call of LikeArtist.
func (mr *MockLikeRepoMockRecorder) LikeArtist(ctx, userID, artistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikeArtist", reflect.TypeOf((*MockLikeRepo)(nil).LikeArtist), ctx, userID, artistID)
}

// LikeTrack mocks base method.
func (m *MockLikeRepo) LikeTrack(ctx context.Context, userID, trackID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikeTrack", ctx, userID, trackID)
	ret0, _ := ret[0].(error)
	return ret0
}

// LikeTrack indicates an expected call of LikeTrack.
func (mr *MockLikeRepoMockRecorder) LikeTrack(ctx, userID, trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikeTrack", reflect.TypeOf((*MockLikeRepo)(nil).LikeTrack), ctx, userID, trackID)
}

// UnlikeArtist mocks base method.
func (m *MockLikeRepo) UnlikeArtist(ctx context.Context, userID, artistID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlikeArtist", ctx, userID, artistID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlikeArtist indicates an expected call of UnlikeArtist.
func (mr *MockLikeRepoMockRecorder) UnlikeArtist(ctx, userID, artistID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlikeArtist", reflect.TypeOf((*MockLikeRepo)(nil).UnlikeArtist), ctx, userID, artistID)
}

// UnlikeTrack mocks base method.
func (m *MockLikeRepo) UnlikeTrack(ctx context.Context, userID, trackID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlikeTrack", ctx, userID, trackID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlikeTrack indicates an expected call of UnlikeTrack.
func (mr *MockLikeRepoMockRecorder) UnlikeTrack(ctx, userID, trackID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlikeTrack", reflect.TypeOf((*MockLikeRepo)(nil).UnlikeTrack), ctx, userID, trackID)
}

// MockTrackClient is a mock of TrackClient interface.
type MockTrackClient struct {
	ctrl     *gomock.Controller
	recorder *MockTrackClientMockRecorder
}

// MockTrackClientMockRecorder is the mock recorder for MockTrackClient.
type MockTrackClientMockRecorder struct {
	mock *MockTrackClient
}

// NewMockTrackClient creates a new mock instance.
func NewMockTrackClient(ctrl *gomock.Controller) *MockTrackClient {
	mock := &MockTrackClient{ctrl: ctrl}
	mock.recorder = &MockTrackClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrackClient) EXPECT() *MockTrackClientMockRecorder {
	return m.recorder
}

// GetTrack mocks base method.
func (m *MockTrackClient) GetTrack(ctx context.Context, id uuid.UUID) (models.Track, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrack", ctx, id)
	ret0, _ := ret[0].(models.Track)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrack indicates an expected call of GetTrack.
func (mr *MockTrackClientMockRecorder) GetTrack(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrack", reflect.TypeOf((*MockTrackClient)(nil).GetTrack), ctx, id)
}
