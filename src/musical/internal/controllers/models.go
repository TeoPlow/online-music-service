// Package controllers содержит основные контроллеры приложения (http и grpc)
package controllers

import (
	"context"

	"github.com/google/uuid"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
)

type MusicService interface {
	AddMusic(context.Context, models.Track) error
	GetMusic(context.Context, uuid.UUID) (models.Track, error)
}
