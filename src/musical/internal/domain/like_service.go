package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models/dto"
	"github.com/TeoPlow/online-music-service/src/musical/internal/storage"
)

type LikeRepo interface {
	LikeTrack(ctx context.Context, userID, trackID uuid.UUID) error
	UnlikeTrack(ctx context.Context, userID, trackID uuid.UUID) error
	LikeArtist(ctx context.Context, userID, artistID uuid.UUID) error
	UnlikeArtist(ctx context.Context, userID, artistID uuid.UUID) error
	GetLikedTracks(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Track, error)
	GetLikedArtists(ctx context.Context, userID uuid.UUID, offset, limit int) ([]models.Artist, error)
}

type TrackClient interface {
	GetTrack(ctx context.Context, id uuid.UUID) (models.Track, error)
}

type LikeService struct {
	likeRepo     LikeRepo
	artistClient ArtistClient
	trackClient  TrackClient
	txm          TxManager
}

func NewLikeService(
	likeRepo LikeRepo,
	artistClient ArtistClient,
	trackClient TrackClient,
	txm TxManager,
) *LikeService {
	return &LikeService{
		likeRepo:     likeRepo,
		artistClient: artistClient,
		trackClient:  trackClient,
		txm:          txm,
	}
}

func (s *LikeService) LikeTrack(ctx context.Context, userID, trackID uuid.UUID) error {
	return s.txm.RunSerializable(ctx, func(ctx context.Context) error {
		_, err := s.trackClient.GetTrack(ctx, trackID)
		if err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return ErrNotFound
			}
			return ErrInternal
		}
		if err := s.likeRepo.LikeTrack(ctx, userID, trackID); err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return ErrNotFound
			}
			return ErrInternal
		}
		SendLikeTrack(ctx, userID, trackID)
		return nil
	})
}

func (s *LikeService) UnlikeTrack(ctx context.Context, userID, trackID uuid.UUID) error {
	return s.txm.RunSerializable(ctx, func(ctx context.Context) error {
		_, err := s.trackClient.GetTrack(ctx, trackID)
		if err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return ErrNotFound
			}
			return ErrInternal
		}
		if err := s.likeRepo.UnlikeTrack(ctx, userID, trackID); err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return ErrNotFound
			}
			return ErrInternal
		}
		return nil
	})
}

func (s *LikeService) LikeArtist(ctx context.Context, userID, artistID uuid.UUID) error {
	return s.txm.RunSerializable(ctx, func(ctx context.Context) error {
		_, err := s.artistClient.GetArtist(ctx, artistID)
		if err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return ErrNotFound
			}
			return ErrInternal
		}
		if err := s.likeRepo.LikeArtist(ctx, userID, artistID); err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return ErrNotFound
			}
			return ErrInternal
		}
		SendLikeArtist(ctx, userID, artistID)
		return nil
	})
}

func (s *LikeService) UnlikeArtist(ctx context.Context, userID, artistID uuid.UUID) error {
	return s.txm.RunSerializable(ctx, func(ctx context.Context) error {
		_, err := s.artistClient.GetArtist(ctx, artistID)
		if err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return ErrNotFound
			}
			return ErrInternal
		}
		if err := s.likeRepo.UnlikeArtist(ctx, userID, artistID); err != nil {
			if errors.Is(err, storage.ErrNotExists) {
				return ErrNotFound
			}
			return ErrInternal
		}
		return nil
	})
}

func (s *LikeService) GetLikedTracks(
	ctx context.Context,
	userID uuid.UUID,
	req dto.GetLikedTracksRequest,
) ([]models.Track, error) {
	offset := (req.Page - 1) * req.PageSize
	artists, err := s.likeRepo.GetLikedTracks(ctx, userID, offset, req.PageSize)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return artists, nil
}

func (s *LikeService) GetLikedArtists(
	ctx context.Context,
	userID uuid.UUID,
	req dto.GetLikedArtistsRequest,
) ([]models.Artist, error) {
	offset := (req.Page - 1) * req.PageSize
	artists, err := s.likeRepo.GetLikedArtists(ctx, userID, offset, req.PageSize)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return artists, nil
}
