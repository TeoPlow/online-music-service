package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models/dto"
	"github.com/TeoPlow/online-music-service/src/musical/internal/storage"
)

type AlbumRepo interface {
	repo[models.Album]

	GetByFilter(ctx context.Context,
		offset, limit int,
		searchQuery *string,
		albumID *uuid.UUID) ([]models.Album, error)
}

type ArtistClient interface {
	GetArtist(ctx context.Context, id uuid.UUID) (models.Artist, error)
}

type AlbumService struct {
	repo         AlbumRepo
	txm          TxManager
	artistClient ArtistClient
}

func NewAlbumService(s AlbumRepo, txm TxManager, client ArtistClient) *AlbumService {
	return &AlbumService{repo: s, txm: txm, artistClient: client}
}

func (service *AlbumService) CreateAlbum(ctx context.Context, req dto.CreateAlbumRequest) (models.Album, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return models.Album{}, fmt.Errorf("failed to create uuid: %w", err)
	}
	album := models.Album{
		ID:          id,
		Title:       req.Title,
		ArtistID:    req.ArtistID,
		ReleaseDate: time.Now(),
	}
	if err := service.txm.RunReadUncommited(ctx, func(ctx context.Context) error {
		hasCopy, err := service.repo.FindCopy(ctx, album)
		if err != nil {
			return ErrInternal
		}
		if hasCopy {
			return fmt.Errorf("album %w", ErrAlreadyExists)
		}

		if _, err = service.artistClient.GetArtist(ctx, album.ArtistID); err != nil {
			return err
		}

		if err = service.repo.Add(ctx, album); err != nil {
			return ErrInternal
		}
		return nil
	}); err != nil {
		return models.Album{}, err
	}
	SendAlbumCreate(ctx, album)
	return album, nil
}

func updateAlbumData(data *models.Album, req dto.UpdateAlbumRequest) {
	if data == nil {
		return
	}

	if req.ArtistID != nil {
		data.ArtistID = *req.ArtistID
	}
	if req.Title != nil {
		data.Title = *req.Title
	}
	if req.ReleaseDate != nil {
		data.ReleaseDate = *req.ReleaseDate
	}
}

func (service *AlbumService) UpdateAlbum(ctx context.Context, req dto.UpdateAlbumRequest) (models.Album, error) {
	var album models.Album
	if err := service.txm.RunSerializable(ctx, func(ctx context.Context) error {
		var err error
		album, err = service.repo.GetByID(ctx, req.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return fmt.Errorf("album %w", ErrNotFound)
			}
			return ErrInternal
		}

		updateAlbumData(&album, req)

		if req.ArtistID != nil {
			if _, err = service.artistClient.GetArtist(ctx, *req.ArtistID); err != nil {
				return err
			}
		}

		if err = service.repo.Update(ctx, album); err != nil {
			return ErrInternal
		}
		return nil
	}); err != nil {
		return models.Album{}, err
	}

	return album, nil
}

func (service *AlbumService) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	if err := service.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			return fmt.Errorf("album %w", ErrNotFound)
		}
		return ErrInternal
	}
	return nil
}

func (service *AlbumService) GetAlbum(ctx context.Context, id uuid.UUID) (models.Album, error) {
	album, err := service.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			return models.Album{}, fmt.Errorf("album %w", ErrNotFound)
		}
		return models.Album{}, ErrInternal
	}
	return album, nil
}

func (service *AlbumService) ListAlbums(ctx context.Context, req dto.ListAlbumsRequest) ([]models.Album, error) {
	offset := (req.Page - 1) * req.PageSize
	albums, err := service.repo.GetByFilter(ctx, offset, req.PageSize, req.SearchQuery, req.ArtistID)
	if err != nil {
		return nil, ErrInternal
	}
	return albums, nil
}
