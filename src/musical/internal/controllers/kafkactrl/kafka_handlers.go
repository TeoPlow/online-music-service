// Package kafkactrl содержит обработку сообщений из kafka consumer
package kafkactrl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/TeoPlow/online-music-service/src/musical/internal/controllers"
	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"

	"github.com/google/uuid"
)

const (
	topicArtistCreated = "artist.created"
	topicArtistUpdated = "artist.updated"
	topicArtistDeleted = "artist.deleted"
)

type ArtistDeletedEvent struct {
	SchemaVersion int    `json:"schema_version"`
	ID            string `json:"id"`
}

type KafkaArtistHandler struct {
	artistService controllers.ArtistService
}

func NewKafkaArtistHandler(
	artistService controllers.ArtistService,
) *KafkaArtistHandler {
	return &KafkaArtistHandler{
		artistService: artistService,
	}
}

func (h *KafkaArtistHandler) HandleMessage(ctx context.Context, topic string, message []byte) error {
	switch topic {
	case topicArtistCreated:
		return h.HandleArtistCreatedEvent(ctx, message)
	case topicArtistUpdated:
		return h.HandleArtistUpdatedEvent(ctx, message)
	case topicArtistDeleted:
		return h.HandleArtistDeletedEvent(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

func (h *KafkaArtistHandler) GetTopics() []string {
	return []string{
		topicArtistCreated,
		topicArtistDeleted,
		topicArtistUpdated,
	}
}

func (h *KafkaArtistHandler) HandleArtistCreatedEvent(ctx context.Context, message []byte) error {
	var event models.Artist
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal artist created event: %w", err)
	}

	if err := h.artistService.CreateArtist(ctx, event); err != nil {
		logger.Logger.Error("failed to create artist",
			"error", err.Error(),
			"where", "KafkaArtistHandler.HandleArtistCreated")
		return err
	}
	logger.Logger.Info("artist created successfully",
		"artist_id", event.ID,
		"where", "KafkaArtistHandler.HandleArtistCreated")
	return nil
}

func (h *KafkaArtistHandler) HandleArtistUpdatedEvent(ctx context.Context, message []byte) error {
	var event models.Artist
	if err := json.Unmarshal(message, &event); err != nil {
		return errors.New("failed to unmarshal artist updated event")
	}

	if err := h.artistService.UpdateArtist(ctx, event); err != nil {
		logger.Logger.Error("failed to update artist",
			slog.String("artist_id", event.ID.String()),
			slog.String("error", err.Error()))
		return err
	}

	logger.Logger.Info("artist updated from Kafka event",
		slog.String("artist_id", event.ID.String()))
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
