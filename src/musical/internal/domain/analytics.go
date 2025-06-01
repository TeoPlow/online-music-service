package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
)

const (
	topicAnalyticArtist     = "music-artists"
	topicAnalyticAlbum      = "music-albums"
	topicAnalyticTrack      = "music-tracks"
	topicAnalyticLikeTrack  = "music-liked-tracks"
	topicAnalyticLikeArtist = "music-liked-artists"
	topicAnalyticListen     = "music-events"
)

type Publisher interface {
	Publish(
		ctx context.Context,
		topic string,
		key string,
		payload []byte,
	) error
}

var analyticPublisher Publisher

func InitAnalyticSender(p Publisher) {
	analyticPublisher = p
}

func SendArtistCreate(ctx context.Context, data models.Artist) {
	l := logger.Logger.With(slog.String("where", "domain.SendArtistCreate"))
	if analyticPublisher == nil {
		l.Error("analytic publisher uninitialized, can not send statistics")
		return
	}
	payload, err := json.Marshal(data)
	if err != nil {
		l.Error("failed to marshal payload",
			slog.String("error", err.Error()))
	}
	if err = analyticPublisher.Publish(ctx,
		topicAnalyticArtist,
		data.ID.String(), payload); err != nil {
		l.Error("failed to send statistics",
			slog.String("error", err.Error()))
	}
}

func SendAlbumCreate(ctx context.Context, data models.Album) {
	l := logger.Logger.With(slog.String("where", "domain.SendAlbumCreate"))
	if analyticPublisher == nil {
		l.Error("analytic publisher uninitialized, can not send statistics")
		return
	}
	payload, err := json.Marshal(data)
	if err != nil {
		l.Error("failed to marshal payload",
			slog.String("error", err.Error()))
	}
	if err = analyticPublisher.Publish(ctx,
		topicAnalyticAlbum,
		data.ID.String(), payload); err != nil {
		l.Error("failed to send statistics",
			slog.String("error", err.Error()))
	}
}

func SendTrackCreate(ctx context.Context, data models.Track) {
	l := logger.Logger.With(slog.String("where", "domain.SendTrackCreate"))
	if analyticPublisher == nil {
		l.Error("analytic publisher uninitialized, can not send statistics")
		return
	}
	payload, err := json.Marshal(data)
	if err != nil {
		l.Error("failed to marshal payload",
			slog.String("error", err.Error()))
	}
	if err = analyticPublisher.Publish(ctx,
		topicAnalyticTrack,
		data.ID.String(), payload); err != nil {
		l.Error("failed to send statistics",
			slog.String("error", err.Error()))
	}
}

func SendLikeTrack(ctx context.Context, userID, trackID uuid.UUID) {
	l := logger.Logger.With(slog.String("where", "domain.SendLikeTrack"))
	if analyticPublisher == nil {
		l.Error("analytic publisher uninitialized, can not send statistics")
		return
	}
	payload := fmt.Sprintf(`{"user_id":"%s", "track_id":"%s"}`, userID, trackID)
	if err := analyticPublisher.Publish(ctx,
		topicAnalyticLikeTrack,
		userID.String(), []byte(payload)); err != nil {
		l.Error("failed to send statistics",
			slog.String("error", err.Error()))
	}
}

func SendLikeArtist(ctx context.Context, userID, artistID uuid.UUID) {
	l := logger.Logger.With(slog.String("where", "domain.SendLikeTrack"))
	if analyticPublisher == nil {
		l.Error("analytic publisher uninitialized, can not send statistics")
		return
	}
	payload := fmt.Sprintf(`{"user_id":"%s", "track_id":"%s"}`, userID, artistID)
	if err := analyticPublisher.Publish(ctx,
		topicAnalyticLikeArtist,
		userID.String(), []byte(payload)); err != nil {
		l.Error("failed to send statistics",
			slog.String("error", err.Error()))
	}
}

func SendListeningEvent(ctx context.Context, userID, trackID uuid.UUID) {
	l := logger.Logger.With(slog.String("where", "domain.SendLikeTrack"))
	if analyticPublisher == nil {
		l.Error("analytic publisher uninitialized, can not send statistics")
		return
	}
	payload := fmt.Sprintf(`{"user_id":"%s", "track_id":"%s"}`, userID, trackID)
	if err := analyticPublisher.Publish(ctx,
		topicAnalyticListen,
		userID.String(), []byte(payload)); err != nil {
		l.Error("failed to send statistics",
			slog.String("error", err.Error()))
	}
}
