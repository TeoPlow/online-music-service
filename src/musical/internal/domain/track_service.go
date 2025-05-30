package domain

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models/dto"
	"github.com/TeoPlow/online-music-service/src/musical/internal/storage"
)

type TrackRepo interface {
	repo[models.Track]

	GetByFilter(ctx context.Context,
		offset, limit int,
		genre, searchQuery *string,
		albumID *uuid.UUID) ([]models.Track, error)
}

type AlbumClient interface {
	GetAlbum(ctx context.Context, id uuid.UUID) (models.Album, error)
}

type StreamingClient interface {
	SaveTrack(context.Context, uuid.UUID, *bytes.Buffer) error
	DeleteTrack(context.Context, uuid.UUID) error
}

type TrackService struct {
	repo            TrackRepo
	txm             TxManager
	albumClient     AlbumClient
	streamingClient StreamingClient
}

func NewTrackService(s TrackRepo, txm TxManager,
	client AlbumClient, streamer StreamingClient,
) *TrackService {
	return &TrackService{
		repo:            s,
		txm:             txm,
		albumClient:     client,
		streamingClient: streamer,
	}
}

func (service *TrackService) CreateTrack(ctx context.Context,
	req dto.CreateTrackRequest, file *bytes.Buffer,
) (models.Track, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return models.Track{}, fmt.Errorf("failed to create uuid: %w", err)
	}
	now := time.Now()
	track := models.Track{
		ID:         id,
		Title:      req.Title,
		AlbumID:    req.AlbumID,
		Genre:      req.Genre,
		Duration:   req.Duration,
		Lyrics:     req.Lyrics,
		IsExplicit: req.IsExplicit,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := service.txm.RunSerializable(ctx, func(ctx context.Context) error {
		hasCopy, err := service.repo.FindCopy(ctx, track)
		if err != nil {
			return ErrInternal
		}
		if hasCopy {
			return fmt.Errorf("track %w", ErrAlreadyExists)
		}

		if _, err = service.albumClient.GetAlbum(ctx, track.AlbumID); err != nil {
			return err
		}

		if err = service.repo.Add(ctx, track); err != nil {
			return ErrInternal
		}

		if err = service.streamingClient.SaveTrack(ctx, track.ID, file); err != nil {
			return ErrInternal
		}
		return nil
	}); err != nil {
		return models.Track{}, err
	}

	return track, nil
}

func updateTrackData(data *models.Track, req dto.UpdateTrackRequest) {
	if data == nil {
		return
	}

	if req.AlbumID != nil {
		data.AlbumID = *req.AlbumID
	}
	if req.Title != nil {
		data.Title = *req.Title
	}
	if req.Genre != nil {
		data.Genre = *req.Genre
	}
	if req.Duration != nil {
		data.Duration = *req.Duration
	}
	if req.Lyrics != nil {
		data.Lyrics = req.Lyrics
	}
	if req.IsExplicit != nil {
		data.IsExplicit = *req.IsExplicit
	}
	data.UpdatedAt = time.Now()
}

func (service *TrackService) UpdateTrack(ctx context.Context, req dto.UpdateTrackRequest) (models.Track, error) {
	var track models.Track
	if err := service.txm.RunSerializable(ctx, func(ctx context.Context) error {
		var err error
		track, err = service.repo.GetByID(ctx, req.ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return fmt.Errorf("track %w", ErrNotFound)
			}
			return ErrInternal
		}

		updateTrackData(&track, req)

		if req.AlbumID != nil {
			if _, err = service.albumClient.GetAlbum(ctx, *req.AlbumID); err != nil {
				return err
			}
		}

		if err = service.repo.Update(ctx, track); err != nil {
			return ErrInternal
		}
		return nil
	}); err != nil {
		return models.Track{}, err
	}

	return track, nil
}

func (service *TrackService) DeleteTrack(ctx context.Context, id uuid.UUID) error {
	return service.txm.RunReadUncommited(ctx, func(ctx context.Context) error {
		if err := service.repo.Delete(ctx, id); err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return fmt.Errorf("track %w", ErrNotFound)
			}
			return ErrInternal
		}
		if err := service.streamingClient.DeleteTrack(ctx, id); err != nil {
			return ErrInternal
		}
		return nil
	})
}

func (service *TrackService) GetTrack(ctx context.Context, id uuid.UUID) (models.Track, error) {
	track, err := service.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			return models.Track{}, fmt.Errorf("track %w", ErrNotFound)
		}
		return models.Track{}, ErrInternal
	}
	return track, nil
}

func (service *TrackService) ListTracks(ctx context.Context, req dto.ListTracksRequest) ([]models.Track, error) {
	offset := (req.Page - 1) * req.PageSize
	tracks, err := service.repo.GetByFilter(ctx,
		offset, req.PageSize,
		req.Genre, req.SearchQuery,
		req.AlbumID)
	if err != nil {
		return nil, ErrInternal
	}
	return tracks, nil
}
