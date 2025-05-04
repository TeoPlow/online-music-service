// Package domain содержит логику запросов (валидации и т.д.) внутри структур Service.
package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
)

type MusicRepo interface {
	Add(ctx context.Context, track models.Track) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, track models.Track) error
	GetByID(ctx context.Context, id uuid.UUID) (models.Track, error)
	GetAll(ctx context.Context) ([]models.Track, error)
}

type TxManager interface {
	RunSerializable(context.Context, func(context.Context) error) error
	RunReadUncommited(context.Context, func(context.Context) error) error
}

type MusicService struct {
	repo MusicRepo
	txm  TxManager
}

func NewMusicService(s MusicRepo, txm TxManager) *MusicService {
	return &MusicService{repo: s, txm: txm}
}

func (service *MusicService) AddMusic(ctx context.Context, track models.Track) error {
	if err := service.repo.Add(ctx, track); err != nil {
		return fmt.Errorf("MusicService.AddMusic: %w", err)
	}
	return nil
}

func (service *MusicService) GetMusic(ctx context.Context, id uuid.UUID) (models.Track, error) {
	res, err := service.repo.GetByID(ctx, id)
	if err != nil {
		return models.Track{}, fmt.Errorf("MusicService.GetMusic: %w", err)
	}
	return res, nil
}
