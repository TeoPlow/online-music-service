package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/TeoPlow/online-music-service/src/musical/internal/domain"
	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"

	"github.com/google/uuid"
)

type ArtistEvent struct {
	SchemaVersion int       `json:"schema_version"`
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Author        string    `json:"author"`
	Producer      string    `json:"producer"`
	Country       string    `json:"country"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ArtistDeletedEvent struct {
	SchemaVersion int    `json:"schema_version"`
	ID            string `json:"id"`
}

type KafkaArtistHandler struct {
	artistService domain.ArtistService
	topics        struct {
		ArtistCreated string `yaml:"artist_created"`
		ArtistUpdated string `yaml:"artist_updated"`
		ArtistDeleted string `yaml:"artist_deleted"`
	}
}

func NewKafkaArtistHandler(
	artistService domain.ArtistService,
	topics struct {
		ArtistCreated string `yaml:"artist_created"`
		ArtistUpdated string `yaml:"artist_updated"`
		ArtistDeleted string `yaml:"artist_deleted"`
	},
) *KafkaArtistHandler {
	return &KafkaArtistHandler{
		artistService: artistService,
		topics:        topics,
	}
}

func (h *KafkaArtistHandler) HandleMessage(ctx context.Context, topic string, message []byte) error {
	switch topic {
	case h.topics.ArtistCreated:
		return h.HandleArtistCreatedEvent(ctx, message)
	case h.topics.ArtistUpdated:
		return h.HandleArtistUpdatedEvent(ctx, message)
	case h.topics.ArtistDeleted:
		return h.HandleArtistDeletedEvent(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

func (h *KafkaArtistHandler) HandleArtistCreatedEvent(ctx context.Context, message []byte) error {
	var event ArtistEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal artist created event: %w", err)
	}
	artistID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid artist ID: %w", err)
	}

	artist := models.Artist{
		ID:          artistID,
		Name:        event.Name,
		Author:      event.Author,
		Producer:    event.Producer,
		Country:     event.Country,
		Description: event.Description,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   time.Now(),
	}
	if err := h.artistService.AddArtist(ctx, artist); err != nil {
		logger.Logger.Error("failed to create artist",
			"error", err.Error(),
			"where", "KafkaArtistHandler.HandleArtistCreated")
		return err
	}
	logger.Logger.Info("artist created successfully",
		"artist_id", artist.ID,
		"where", "KafkaArtistHandler.HandleArtistCreated")
	return nil

}

func (h *KafkaArtistHandler) HandleArtistUpdatedEvent(ctx context.Context, message []byte) error {
	var event ArtistEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return errors.New("failed to unmarshal artist updated event")
	}

	artistID, err := uuid.Parse(event.ID)
	if err != nil {
		return errors.New("invalid artist ID in event")
	}
	existingArtist, err := h.artistService.GetArtist(ctx, artistID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return h.HandleArtistCreatedEvent(ctx, message)
		}
		return err
	}

	updatedArtist := models.Artist{
		ID:          artistID,
		Name:        event.Name,
		Author:      event.Author,
		Producer:    event.Producer,
		Country:     event.Country,
		Description: event.Description,
		CreatedAt:   existingArtist.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	if err := h.artistService.UpdateArtist(ctx, updatedArtist); err != nil {
		logger.Logger.Error("failed to update artist",
			slog.String("artist_id", artistID.String()),
			slog.String("error", err.Error()))
		return err
	}

	logger.Logger.Info("artist updated from Kafka event",
		slog.String("artist_id", artistID.String()))
	return nil
}

func (h *KafkaArtistHandler) HandleArtistDeletedEvent(ctx context.Context, message []byte) error {
	var event ArtistDeletedEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal artist deleted event: %w", err)
	}

	artistID, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("invalid artist ID: %w", err)
	}

	if err := h.artistService.DeleteArtist(ctx, artistID); err != nil {
		logger.Logger.Error("failed to delete artist",
			slog.String("artist_id", artistID.String()),
			slog.String("error", err.Error()))
		return err
	}

	logger.Logger.Info("artist deleted successfully",
		slog.String("artist_id", artistID.String()))
	return nil
}
