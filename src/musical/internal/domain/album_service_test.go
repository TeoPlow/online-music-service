package domain_test

import (
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

func TestCreateAlbum(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	model := models.Album{
		ID:          id,
		Title:       "album title",
		ArtistID:    uuid.New(),
		ReleaseDate: time.Now(),
	}
	req := dto.CreateAlbumRequest{
		Title:    "album title",
		ArtistID: model.ArtistID,
	}
	tests := []struct {
		name        string
		repoFindRes bool
		repoAddErr  error
		clientErr   error
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
			name:      "no artist",
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
	}

	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockAlbumRepo(ctrl)
	mockTxm := mocks.NewMockTxManager(ctrl)
	mockClient := mocks.NewMockArtistClient(ctrl)
	service := domain.NewAlbumService(
		mockRepo,
		mockTxm,
		mockClient,
	)

	mockTxm.EXPECT().RunReadUncommited(
		gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(context.Context) error) error {
			return f(ctx)
		}).AnyTimes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().FindCopy(gomock.Any(), gomock.Any()).Return(tt.repoFindRes, nil)
			if !tt.repoFindRes {
				if tt.clientErr == nil {
					mockRepo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(tt.repoAddErr)
				}
				mockClient.EXPECT().GetArtist(gomock.Any(), model.ArtistID).Return(models.Artist{}, tt.clientErr)
			}

			res, err := service.CreateAlbum(t.Context(), req)

			assert.ErrorIs(t, err, tt.wantErr)

			if err == nil {
				assert.NotEqual(t, uuid.UUID{}, res.ID)
				assert.Equal(t, req.ArtistID, res.ArtistID)
				assert.Equal(t, req.Title, res.Title)
				assert.NotEqual(t, time.Time{}, res.ReleaseDate)
			}
		})
	}
}

func TestUpdateAlbum(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	model := models.Album{
		ID:          id,
		Title:       "album title",
		ArtistID:    uuid.New(),
		ReleaseDate: time.Now(),
	}
	fullUpdate := dto.UpdateAlbumRequest{
		ID:          id,
		Title:       testutils.ToPtr("other title"),
		ArtistID:    testutils.ToPtr(uuid.New()),
		ReleaseDate: testutils.ToPtr(time.Now()),
	}
	tests := []struct {
		name          string
		repoGetErr    error
		repoUpdateErr error
		clientErr     error
		old           models.Album
		req           dto.UpdateAlbumRequest
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
			req: dto.UpdateAlbumRequest{
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
			name:      "no artist",
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
	mockRepo := mocks.NewMockAlbumRepo(ctrl)
	mockTxm := mocks.NewMockTxManager(ctrl)
	mockClient := mocks.NewMockArtistClient(ctrl)
	service := domain.NewAlbumService(
		mockRepo,
		mockTxm,
		mockClient,
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
				if tt.req.ArtistID != nil {
					mockClient.EXPECT().GetArtist(gomock.Any(), *tt.req.ArtistID).Return(models.Artist{}, tt.clientErr)
				}
			}

			res, err := service.UpdateAlbum(t.Context(), tt.req)

			assert.ErrorIs(t, err, tt.wantErr)

			if err == nil {
				assert.Equal(t, tt.req.ID, res.ID)
				if tt.req.Title != nil {
					assert.Equal(t, *tt.req.Title, res.Title)
				} else {
					assert.Equal(t, tt.old.Title, res.Title)
				}
				if tt.req.ArtistID != nil {
					assert.Equal(t, *tt.req.ArtistID, res.ArtistID)
				} else {
					assert.Equal(t, tt.old.ArtistID, res.ArtistID)
				}
				if tt.req.ReleaseDate != nil {
					assert.Equal(t, *tt.req.ReleaseDate, res.ReleaseDate)
				} else {
					assert.Equal(t, tt.old.ReleaseDate, res.ReleaseDate)
				}
			}
		})
	}
}

func TestGetAlbum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		repoErr   error
		repoModel models.Album
	}{
		{
			name:    "success",
			repoErr: nil,
			repoModel: models.Album{
				ID:          uuid.New(),
				Title:       "album title",
				ArtistID:    uuid.New(),
				ReleaseDate: time.Now(),
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
	mockRepo := mocks.NewMockAlbumRepo(ctrl)
	service := domain.NewAlbumService(
		mockRepo,
		mocks.NewMockTxManager(ctrl),
		mocks.NewMockArtistClient(ctrl),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().GetByID(gomock.Any(), tt.repoModel.ID).Return(tt.repoModel, tt.repoErr)

			res, err := service.GetAlbum(t.Context(), tt.repoModel.ID)
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

func TestDeleteAlbum(t *testing.T) {
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
	mockRepo := mocks.NewMockAlbumRepo(ctrl)
	service := domain.NewAlbumService(
		mockRepo,
		mocks.NewMockTxManager(ctrl),
		mocks.NewMockArtistClient(ctrl),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := uuid.New()
			mockRepo.EXPECT().Delete(gomock.Any(), id).Return(tt.repoErr)

			err := service.DeleteAlbum(t.Context(), id)
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
