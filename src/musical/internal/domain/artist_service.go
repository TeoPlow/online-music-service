package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	"github.com/TeoPlow/online-music-service/src/musical/internal/storage"
)

type ArtistRepo repo[models.Artist]

type ArtistService struct {
	repo ArtistRepo
	txm  TxManager
}

func NewArtistService(s ArtistRepo, txm TxManager) *ArtistService {
	return &ArtistService{repo: s, txm: txm}
}

func (service *ArtistService) GetArtist(ctx context.Context, id uuid.UUID) (models.Artist, error) {
	artist, err := service.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			return models.Artist{}, fmt.Errorf("artist %w", ErrNotFound)
		}
		return models.Artist{}, ErrInternal
	}
	return artist, nil
}
