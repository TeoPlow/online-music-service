package domain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TeoPlow/online-music-service/src/musical/internal/domain"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models/dto"
	"github.com/TeoPlow/online-music-service/src/musical/internal/storage"
	"github.com/TeoPlow/online-music-service/src/musical/tests/mocks"
)

func TestLikeTrack(t *testing.T) {
	t.Parallel()

	trackID := uuid.New()
	userID := uuid.New()
	track := models.Track{
		ID:        trackID,
		Title:     "Test Track",
		Duration:  180 * time.Second,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name        string
		setupMocks  func(*mocks.MockLikeRepo, *mocks.MockTrackClient, *mocks.MockTxManager)
		expectedErr error
	}{
		{
			name: "success",
			setupMocks: func(likeRepo *mocks.MockLikeRepo, trackRepo *mocks.MockTrackClient, txm *mocks.MockTxManager) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				trackRepo.EXPECT().
					GetTrack(gomock.Any(), trackID).
					Return(track, nil)

				likeRepo.EXPECT().
					LikeTrack(gomock.Any(), userID, trackID).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "track not exists",
			setupMocks: func(likeRepo *mocks.MockLikeRepo, trackRepo *mocks.MockTrackClient, txm *mocks.MockTxManager) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				trackRepo.EXPECT().
					GetTrack(gomock.Any(), trackID).
					Return(models.Track{}, storage.ErrNotExists)
			},
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "track repo error",
			setupMocks: func(likeRepo *mocks.MockLikeRepo, trackRepo *mocks.MockTrackClient, txm *mocks.MockTxManager) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				trackRepo.EXPECT().
					GetTrack(gomock.Any(), trackID).
					Return(models.Track{}, errors.New("internal error"))
			},
			expectedErr: domain.ErrInternal,
		},
		{
			name: "like repo error",
			setupMocks: func(likeRepo *mocks.MockLikeRepo, trackRepo *mocks.MockTrackClient, txm *mocks.MockTxManager) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				trackRepo.EXPECT().
					GetTrack(gomock.Any(), trackID).
					Return(track, nil)

				likeRepo.EXPECT().
					LikeTrack(gomock.Any(), userID, trackID).
					Return(errors.New("like error"))
			},
			expectedErr: domain.ErrInternal,
		},
		{
			name: "transaction error",
			setupMocks: func(likeRepo *mocks.MockLikeRepo, trackRepo *mocks.MockTrackClient, txm *mocks.MockTxManager) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					Return(errors.New("transaction error"))
			},
			expectedErr: errors.New("transaction error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeRepo := mocks.NewMockLikeRepo(ctrl)
			trackRepo := mocks.NewMockTrackClient(ctrl)
			artistRepo := mocks.NewMockArtistClient(ctrl)
			txm := mocks.NewMockTxManager(ctrl)

			tt.setupMocks(likeRepo, trackRepo, txm)

			service := domain.NewLikeService(likeRepo, artistRepo, trackRepo, txm)
			err := service.LikeTrack(context.Background(), userID, trackID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, tt.expectedErr) && err.Error() != tt.expectedErr.Error() {
					t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnlikeTrack(t *testing.T) {
	t.Parallel()

	trackID := uuid.New()
	userID := uuid.New()
	track := models.Track{
		ID:        trackID,
		Title:     "Test Track",
		Duration:  180 * time.Second,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name        string
		setupMocks  func(*mocks.MockLikeRepo, *mocks.MockTrackClient, *mocks.MockTxManager)
		expectedErr error
	}{
		{
			name: "success",
			setupMocks: func(likeRepo *mocks.MockLikeRepo, trackRepo *mocks.MockTrackClient, txm *mocks.MockTxManager) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				trackRepo.EXPECT().
					GetTrack(gomock.Any(), trackID).
					Return(track, nil)

				likeRepo.EXPECT().
					UnlikeTrack(gomock.Any(), userID, trackID).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "track not exists",
			setupMocks: func(likeRepo *mocks.MockLikeRepo, trackRepo *mocks.MockTrackClient, txm *mocks.MockTxManager) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				trackRepo.EXPECT().
					GetTrack(gomock.Any(), trackID).
					Return(models.Track{}, storage.ErrNotExists)
			},
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "unlike error",
			setupMocks: func(likeRepo *mocks.MockLikeRepo, trackRepo *mocks.MockTrackClient, txm *mocks.MockTxManager) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				trackRepo.EXPECT().
					GetTrack(gomock.Any(), trackID).
					Return(track, nil)

				likeRepo.EXPECT().
					UnlikeTrack(gomock.Any(), userID, trackID).
					Return(errors.New("unlike error"))
			},
			expectedErr: domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeRepo := mocks.NewMockLikeRepo(ctrl)
			trackRepo := mocks.NewMockTrackClient(ctrl)
			artistRepo := mocks.NewMockArtistClient(ctrl)
			txm := mocks.NewMockTxManager(ctrl)

			tt.setupMocks(likeRepo, trackRepo, txm)

			service := domain.NewLikeService(likeRepo, artistRepo, trackRepo, txm)
			err := service.UnlikeTrack(context.Background(), userID, trackID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, tt.expectedErr) && err.Error() != tt.expectedErr.Error() {
					t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLikeArtist(t *testing.T) {
	t.Parallel()

	artistID := uuid.New()
	userID := uuid.New()
	artist := models.Artist{
		ID:        artistID,
		Name:      "Test Artist",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		setupMocks func(*mocks.MockLikeRepo,
			*mocks.MockArtistClient, *mocks.MockTxManager)
		expectedErr error
	}{
		{
			name: "success",
			setupMocks: func(likeRepo *mocks.MockLikeRepo,
				artistRepo *mocks.MockArtistClient,
				txm *mocks.MockTxManager,
			) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context,
						fn func(context.Context) error,
					) error {
						return fn(ctx)
					})

				artistRepo.EXPECT().
					GetArtist(gomock.Any(), artistID).
					Return(artist, nil)

				likeRepo.EXPECT().
					LikeArtist(gomock.Any(), userID, artistID).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "artist not exists",
			setupMocks: func(likeRepo *mocks.MockLikeRepo,
				artistRepo *mocks.MockArtistClient, txm *mocks.MockTxManager,
			) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context,
						fn func(context.Context) error,
					) error {
						return fn(ctx)
					})

				artistRepo.EXPECT().
					GetArtist(gomock.Any(), artistID).
					Return(models.Artist{}, storage.ErrNotExists)
			},
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeRepo := mocks.NewMockLikeRepo(ctrl)
			trackRepo := mocks.NewMockTrackClient(ctrl)
			artistRepo := mocks.NewMockArtistClient(ctrl)
			txm := mocks.NewMockTxManager(ctrl)

			tt.setupMocks(likeRepo, artistRepo, txm)

			service := domain.NewLikeService(likeRepo, artistRepo, trackRepo, txm)
			err := service.LikeArtist(context.Background(), userID, artistID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, tt.expectedErr) && err.Error() != tt.expectedErr.Error() {
					t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetLikedTracks(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	tracks := []models.Track{
		{
			ID:        uuid.New(),
			Title:     "Track 1",
			Duration:  180 * time.Second,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Title:     "Track 2",
			Duration:  240 * time.Second,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	tests := []struct {
		name        string
		req         dto.GetLikedTracksRequest
		setupMocks  func(*mocks.MockLikeRepo)
		expected    []models.Track
		expectedErr error
	}{
		{
			name: "success",
			req: dto.GetLikedTracksRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMocks: func(likeRepo *mocks.MockLikeRepo) {
				likeRepo.EXPECT().
					GetLikedTracks(gomock.Any(), userID, 0, 10).
					Return(tracks, nil)
			},
			expected:    tracks,
			expectedErr: nil,
		},
		{
			name: "empty result",
			req: dto.GetLikedTracksRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMocks: func(likeRepo *mocks.MockLikeRepo) {
				likeRepo.EXPECT().
					GetLikedTracks(gomock.Any(), userID, 0, 10).
					Return([]models.Track{}, nil)
			},
			expected:    []models.Track{},
			expectedErr: nil,
		},
		{
			name: "repository error",
			req: dto.GetLikedTracksRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMocks: func(likeRepo *mocks.MockLikeRepo) {
				likeRepo.EXPECT().
					GetLikedTracks(gomock.Any(), userID, 0, 10).
					Return(nil, storage.ErrSQL)
			},
			expected:    nil,
			expectedErr: domain.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeRepo := mocks.NewMockLikeRepo(ctrl)
			trackRepo := mocks.NewMockTrackClient(ctrl)
			artistRepo := mocks.NewMockArtistClient(ctrl)
			txm := mocks.NewMockTxManager(ctrl)

			tt.setupMocks(likeRepo)

			service := domain.NewLikeService(likeRepo, artistRepo, trackRepo, txm)
			result, err := service.GetLikedTracks(context.Background(), userID, tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetLikedArtists(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	artists := []models.Artist{
		{
			ID:        uuid.New(),
			Name:      "Artist 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Artist 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	tests := []struct {
		name        string
		req         dto.GetLikedArtistsRequest
		setupMocks  func(*mocks.MockLikeRepo)
		expected    []models.Artist
		expectedErr error
	}{
		{
			name: "success",
			req: dto.GetLikedArtistsRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMocks: func(likeRepo *mocks.MockLikeRepo) {
				likeRepo.EXPECT().
					GetLikedArtists(gomock.Any(), userID, 0, 10).
					Return(artists, nil)
			},
			expected:    artists,
			expectedErr: nil,
		},
		{
			name: "pagination",
			req: dto.GetLikedArtistsRequest{
				Page:     2,
				PageSize: 5,
			},
			setupMocks: func(likeRepo *mocks.MockLikeRepo) {
				likeRepo.EXPECT().
					GetLikedArtists(gomock.Any(), userID, 5, 5).
					Return(artists[1:], nil)
			},
			expected:    artists[1:],
			expectedErr: nil,
		},
		{
			name: "repository error",
			req: dto.GetLikedArtistsRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMocks: func(likeRepo *mocks.MockLikeRepo) {
				likeRepo.EXPECT().
					GetLikedArtists(gomock.Any(), userID, 0, 10).
					Return(nil, storage.ErrNotExists)
			},
			expected:    nil,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeRepo := mocks.NewMockLikeRepo(ctrl)
			trackRepo := mocks.NewMockTrackClient(ctrl)
			artistRepo := mocks.NewMockArtistClient(ctrl)
			txm := mocks.NewMockTxManager(ctrl)

			tt.setupMocks(likeRepo)

			service := domain.NewLikeService(likeRepo, artistRepo, trackRepo, txm)
			result, err := service.GetLikedArtists(context.Background(), userID, tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedErr))
			} else {
				assert.NoError(t, err)
				require.Equal(t, len(tt.expected), len(result))
				for i := range tt.expected {
					assert.Equal(t, tt.expected[i].ID, result[i].ID)
					assert.Equal(t, tt.expected[i].Name, result[i].Name)
				}
			}
		})
	}
}

func TestUnlikeArtist(t *testing.T) {
	t.Parallel()

	artistID := uuid.New()
	userID := uuid.New()
	artist := models.Artist{
		ID:        artistID,
		Name:      "Test Artist",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		setupMocks func(*mocks.MockLikeRepo,
			*mocks.MockArtistClient, *mocks.MockTxManager)
		expectedErr error
	}{
		{
			name: "success",
			setupMocks: func(likeRepo *mocks.MockLikeRepo,
				artistRepo *mocks.MockArtistClient, txm *mocks.MockTxManager,
			) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})

				artistRepo.EXPECT().
					GetArtist(gomock.Any(), artistID).
					Return(artist, nil)

				likeRepo.EXPECT().
					UnlikeArtist(gomock.Any(), userID, artistID).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "artist not exists",
			setupMocks: func(likeRepo *mocks.MockLikeRepo,
				artistRepo *mocks.MockArtistClient,
				txm *mocks.MockTxManager,
			) {
				txm.EXPECT().
					RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context,
						fn func(context.Context) error,
					) error {
						return fn(ctx)
					})

				artistRepo.EXPECT().
					GetArtist(gomock.Any(), artistID).
					Return(models.Artist{}, storage.ErrNotExists)
			},
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			likeRepo := mocks.NewMockLikeRepo(ctrl)
			trackRepo := mocks.NewMockTrackClient(ctrl)
			artistRepo := mocks.NewMockArtistClient(ctrl)
			txm := mocks.NewMockTxManager(ctrl)

			tt.setupMocks(likeRepo, artistRepo, txm)

			service := domain.NewLikeService(likeRepo, artistRepo, trackRepo, txm)
			err := service.UnlikeArtist(context.Background(), userID, artistID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				if !errors.Is(err, tt.expectedErr) && err.Error() != tt.expectedErr.Error() {
					t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
