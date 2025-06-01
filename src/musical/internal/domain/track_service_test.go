package domain_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/TeoPlow/online-music-service/src/musical/internal/domain"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models/dto"
	"github.com/TeoPlow/online-music-service/src/musical/internal/storage"
	"github.com/TeoPlow/online-music-service/src/musical/tests/mocks"
	"github.com/TeoPlow/online-music-service/src/musical/tests/testutils"
)

func TestCreateTrack(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	model := models.Track{
		ID:         id,
		Title:      "track title",
		AlbumID:    uuid.New(),
		Genre:      "rock",
		Duration:   150 * time.Second,
		IsExplicit: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	req := dto.CreateTrackRequest{
		Title:      "track title",
		AlbumID:    model.AlbumID,
		Genre:      model.Genre,
		Duration:   model.Duration,
		IsExplicit: model.IsExplicit,
	}
	tests := []struct {
		name        string
		repoFindRes bool
		repoAddErr  error
		clientErr   error
		streamerErr error
		wantErr     error
	}{
		{
			name: "succes",
		},
		{
			name:        "has copy",
			repoFindRes: true,
			wantErr:     domain.ErrAlreadyExists,
		},
		{
			name:      "no album",
			clientErr: domain.ErrNotFound,
			wantErr:   domain.ErrNotFound,
		},
		{
			name:      "client error",
			clientErr: domain.ErrInternal,
			wantErr:   domain.ErrInternal,
		},
		{
			name:       "repo error",
			repoAddErr: storage.ErrSQL,
			wantErr:    domain.ErrInternal,
		},
		{
			name:        "minio error",
			streamerErr: domain.ErrInternal,
			wantErr:     domain.ErrInternal,
		},
	}

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTrackRepo(ctrl)
	mockTxm := mocks.NewMockTxManager(ctrl)
	mockClient := mocks.NewMockAlbumClient(ctrl)
	mockStreamer := mocks.NewMockStreamingClient(ctrl)
	service := domain.NewTrackService(
		mockRepo,
		mockTxm,
		mockClient,
		mockStreamer,
	)

	mockTxm.EXPECT().RunSerializable(
		gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(context.Context) error) error {
			return f(ctx)
		}).AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().FindCopy(gomock.Any(), gomock.Any()).Return(tt.repoFindRes, nil)
			if !tt.repoFindRes {
				mockClient.EXPECT().GetAlbum(gomock.Any(), model.AlbumID).Return(models.Album{}, tt.clientErr)
				if tt.clientErr == nil {
					mockRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(tt.repoAddErr)
					if tt.repoAddErr == nil {
						mockStreamer.EXPECT().SaveTrack(gomock.Any(), gomock.Any(),
							gomock.Any()).DoAndReturn(func(arg1, arg2 any, buf *bytes.Buffer) error {
							assert.Equal(t, "meow", buf.String())
							return tt.streamerErr
						})
					}
				}
			}

			res, err := service.CreateTrack(t.Context(), req, bytes.NewBuffer([]byte("meow")))

			assert.ErrorIs(t, err, tt.wantErr)

			if err == nil {
				assert.NotEqual(t, uuid.UUID{}, res.ID)
				assert.Equal(t, req.Title, res.Title)
				assert.Equal(t, req.AlbumID, res.AlbumID)
				assert.Equal(t, req.Genre, res.Genre)
				assert.Equal(t, req.Duration, res.Duration)
				assert.Equal(t, req.Lyrics, res.Lyrics)
				assert.Equal(t, req.IsExplicit, res.IsExplicit)
				assert.Equal(t, res.CreatedAt, res.UpdatedAt)
				assert.Equal(t, testutils.DateOnly(model.CreatedAt),
					testutils.DateOnly(res.CreatedAt))
			}
		})
	}
}

func TestUpdateTrack(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	model := models.Track{
		ID:         id,
		Title:      "track title",
		AlbumID:    uuid.New(),
		Genre:      "rock",
		Duration:   150 * time.Second,
		IsExplicit: true,
		CreatedAt:  time.Date(2025, time.May, 23, 12, 0, 0, 0, time.UTC),
		UpdatedAt:  time.Date(2025, time.May, 23, 12, 0, 0, 0, time.UTC),
	}
	fullUpdate := dto.UpdateTrackRequest{
		ID:         id,
		Title:      testutils.ToPtr("other title"),
		AlbumID:    testutils.ToPtr(uuid.New()),
		Genre:      testutils.ToPtr("classic"),
		IsExplicit: testutils.ToPtr(false),
		Lyrics:     testutils.ToPtr("Meow, meow, meow..."),
	}
	tests := []struct {
		name          string
		repoGetErr    error
		repoUpdateErr error
		clientErr     error
		old           models.Track
		req           dto.UpdateTrackRequest
		wantErr       error
	}{
		{
			name: "full update",
			old:  model,
			req:  fullUpdate,
		},
		{
			name: "partial update",
			old:  model,
			req: dto.UpdateTrackRequest{
				ID:    id,
				Title: testutils.ToPtr("other title"),
			},
		},
		{
			name:       "not exists",
			repoGetErr: storage.ErrNotExists,
			old:        model,
			req:        fullUpdate,
			wantErr:    domain.ErrNotFound,
		},
		{
			name:      "no album",
			clientErr: domain.ErrNotFound,
			old:       model,
			req:       fullUpdate,
			wantErr:   domain.ErrNotFound,
		},
		{
			name:      "client error",
			clientErr: domain.ErrInternal,
			old:       model,
			req:       fullUpdate,
			wantErr:   domain.ErrInternal,
		},
		{
			name:          "repo error",
			repoUpdateErr: storage.ErrSQL,
			req:           fullUpdate,
			wantErr:       domain.ErrInternal,
		},
	}

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTrackRepo(ctrl)
	mockTxm := mocks.NewMockTxManager(ctrl)
	mockClient := mocks.NewMockAlbumClient(ctrl)
	service := domain.NewTrackService(
		mockRepo,
		mockTxm,
		mockClient,
		mocks.NewMockStreamingClient(ctrl),
	)

	mockTxm.EXPECT().RunSerializable(
		gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(context.Context) error) error {
			return f(ctx)
		}).AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().GetByID(gomock.Any(), tt.req.ID).Return(tt.old, tt.repoGetErr)
			if tt.repoGetErr == nil {
				if tt.clientErr == nil {
					mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tt.repoUpdateErr)
				}
				if tt.req.AlbumID != nil {
					mockClient.EXPECT().GetAlbum(gomock.Any(), *tt.req.AlbumID).Return(models.Album{}, tt.clientErr)
				}
			}

			res, err := service.UpdateTrack(t.Context(), tt.req)

			assert.ErrorIs(t, err, tt.wantErr)

			if err == nil {
				assert.Equal(t, tt.req.ID, res.ID)
				if tt.req.Title != nil {
					assert.Equal(t, *tt.req.Title, res.Title)
				} else {
					assert.Equal(t, tt.old.Title, res.Title)
				}
				if tt.req.AlbumID != nil {
					assert.Equal(t, *tt.req.AlbumID, res.AlbumID)
				} else {
					assert.Equal(t, tt.old.AlbumID, res.AlbumID)
				}
				if tt.req.Genre != nil {
					assert.Equal(t, *tt.req.Genre, res.Genre)
				} else {
					assert.Equal(t, tt.old.Genre, res.Genre)
				}
				if tt.req.Duration != nil {
					assert.Equal(t, *tt.req.Duration, res.Duration)
				} else {
					assert.Equal(t, tt.old.Duration, res.Duration)
				}
				if tt.req.Lyrics != nil {
					assert.Equal(t, *tt.req.Lyrics, *res.Lyrics)
				} else if tt.old.Lyrics != nil {
					assert.Equal(t, tt.old.Lyrics, res.Lyrics)
				}
				if tt.req.IsExplicit != nil {
					assert.Equal(t, *tt.req.IsExplicit, res.IsExplicit)
				} else {
					assert.Equal(t, tt.old.IsExplicit, res.IsExplicit)
				}
				assert.True(t, tt.old.UpdatedAt.Before(res.UpdatedAt))
			}
		})
	}
}

func TestGetTrack(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		repoErr   error
		repoModel models.Track
	}{
		{
			name:    "success",
			repoErr: nil,
			repoModel: models.Track{
				ID:         uuid.New(),
				Title:      "track title",
				AlbumID:    uuid.New(),
				Genre:      "rock",
				Duration:   150 * time.Second,
				IsExplicit: true,
				CreatedAt:  time.Date(2025, time.May, 23, 12, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, time.May, 23, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name:    "not exists",
			repoErr: storage.ErrNotExists,
		},
		{
			name:    "internal",
			repoErr: storage.ErrSQL,
		},
	}

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTrackRepo(ctrl)
	service := domain.NewTrackService(
		mockRepo,
		mocks.NewMockTxManager(ctrl),
		mocks.NewMockAlbumClient(ctrl),
		mocks.NewMockStreamingClient(ctrl),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().GetByID(gomock.Any(), tt.repoModel.ID).Return(tt.repoModel, tt.repoErr)

			res, err := service.GetTrack(t.Context(), tt.repoModel.ID)
			if err != nil {
				switch tt.repoErr {
				case storage.ErrNotExists:
					assert.ErrorIs(t, err, domain.ErrNotFound)
				case storage.ErrSQL:
					assert.ErrorIs(t, err, domain.ErrInternal)
				default:
					assert.NoError(t, err)
				}
				return
			}

			assert.Equal(t, tt.repoModel, res)
		})
	}
}

func TestDeleteTrack(t *testing.T) {
	tests := []struct {
		name    string
		repoErr error
	}{
		{
			name:    "success",
			repoErr: nil,
		},
		{
			name:    "not exists",
			repoErr: storage.ErrNotExists,
		},
		{
			name:    "internal",
			repoErr: storage.ErrSQL,
		},
	}

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTrackRepo(ctrl)
	mockTxm := mocks.NewMockTxManager(ctrl)
	mockStreamer := mocks.NewMockStreamingClient(ctrl)
	service := domain.NewTrackService(
		mockRepo,
		mockTxm,
		mocks.NewMockAlbumClient(ctrl),
		mockStreamer,
	)

	mockTxm.EXPECT().RunReadUncommited(
		gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(context.Context) error) error {
			return f(ctx)
		}).AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := uuid.New()
			mockRepo.EXPECT().Delete(gomock.Any(), id).Return(tt.repoErr)
			if tt.repoErr == nil {
				mockStreamer.EXPECT().DeleteTrack(gomock.Any(), id)
			}

			err := service.DeleteTrack(t.Context(), id)
			switch tt.repoErr {
			case storage.ErrNotExists:
				assert.ErrorIs(t, err, domain.ErrNotFound)
			case storage.ErrSQL:
				assert.ErrorIs(t, err, domain.ErrInternal)
			default:
				assert.NoError(t, err)
			}
		})
	}
}
