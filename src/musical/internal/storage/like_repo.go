package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
)

type LikeRepository struct {
	repository
}

func NewLikeRepo(db Executor) *LikeRepository {
	return &LikeRepository{
		repository: repository{
			db: db,
		},
	}
}

func (r *LikeRepository) LikeTrack(ctx context.Context,
	userID uuid.UUID,
	trackID uuid.UUID,
) error {
	var exists bool
	err := r.getExecutor(ctx).Get(ctx, &exists, `
        SELECT EXISTS(
            SELECT 1 FROM liked_tracks
            WHERE user_id = $1 AND track_id = $2
        )
    `, userID, trackID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to check track like existence",
			slog.String("error", err.Error()),
			slog.String("where", "LikeRepository.LikeTrack"))
		return ErrSQL
	}

	if exists {
		return nil
	}
	_, err = r.getExecutor(ctx).Exec(ctx, `
	INSERT INTO liked_tracks (user_id, track_id, created_at)
	VALUES ($1, $2, NOW())
	`, userID, trackID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to like track",
			slog.String("error", err.Error()),
			slog.String("where", "LikeRepository.LikeTrack"))
		return ErrSQL
	}
	return nil
}

func (r *LikeRepository) UnlikeTrack(ctx context.Context,
	userID uuid.UUID,
	trackID uuid.UUID,
) error {
	result, err := r.getExecutor(ctx).Exec(ctx, `
	DELETE FROM liked_tracks
	WHERE user_id = $1 AND track_id = $2`, userID, trackID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to unlike track",
			slog.String("error", err.Error()),
			slog.String("where", "LikeRepository.UnlikeTrack"))
		return ErrSQL
	}
	if result.RowsAffected() == 0 {
		return ErrNotExists
	}
	return nil
}

func (r *LikeRepository) LikeArtist(ctx context.Context,
	userID uuid.UUID,
	artistID uuid.UUID,
) error {
	var exists bool
	err := r.getExecutor(ctx).Get(ctx, &exists, `
        SELECT EXISTS(
            SELECT 1 FROM liked_artists
            WHERE user_id = $1 AND artist_id = $2
        )
    `, userID, artistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to check artist like existence",
			slog.String("error", err.Error()),
			slog.String("where", "LikeRepository.LikeArtist"))
		return ErrSQL
	}
	if exists {
		return nil
	}
	_, err = r.getExecutor(ctx).Exec(ctx, `
		INSERT INTO liked_artists (user_id, artist_id, created_at)
		VALUES ($1, $2, NOW())
	`, userID, artistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to like artist",
			slog.String("error", err.Error()),
			slog.String("where", "LikeRepository.LikeArtist"))
		return ErrSQL
	}
	return nil
}
func (r *LikeRepository) UnlikeArtist(ctx context.Context,
	userID uuid.UUID,
	artistID uuid.UUID,
) error {
	result, err := r.getExecutor(ctx).Exec(ctx, `
		DELETE FROM liked_artists
		WHERE user_id = $1 AND artist_id = $2`, userID, artistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to unlike artist",
			slog.String("error", err.Error()),
			slog.String("where", "LikeRepository.UnlikeArtist"))
		return ErrSQL
	}
	if result.RowsAffected() == 0 {
		return ErrNotExists
	}
	return nil
}

func (repo *LikeRepository) GetLikedTracks(ctx context.Context,
	userID uuid.UUID,
	offset,
	limit int,
) ([]models.Track, error) {
	var tracks []trackForDB
	err := repo.getExecutor(ctx).Select(ctx, &tracks, `
        SELECT t.*
        FROM tracks t
        JOIN liked_tracks lt ON t.id = lt.track_id
        WHERE lt.user_id = $1
        ORDER BY lt.created_at DESC
        OFFSET $2 LIMIT $3
    `, userID, offset, limit)

	if err != nil {

		logger.Logger.Error("failed to get liked tracks",
			slog.String("error", err.Error()),
			slog.String("where", "LikeRepository.GetLikedTracks"))
		return nil, ErrSQL
	}

	result := make([]models.Track, 0, len(tracks))
	for _, t := range tracks {
		result = append(result, convertToModel(t))
	}

	return result, nil
}

func (repo *LikeRepository) GetLikedArtists(ctx context.Context,
	userID uuid.UUID,
	offset,
	limit int,
) ([]models.Artist, error) {
	var artists []models.Artist
	err := repo.getExecutor(ctx).Select(ctx, &artists, `
        SELECT a.*
        FROM artists a
        JOIN liked_artists la ON a.id = la.artist_id
        WHERE la.user_id = $1
        ORDER BY la.created_at DESC
        OFFSET $2 LIMIT $3
    `, userID, offset, limit)

	if err != nil {

		logger.Logger.Error("failed to get liked artists",
			slog.String("error", err.Error()),
			slog.String("where", "LikeRepository.GetLikedArtists"))
		return nil, ErrSQL
	}

	return artists, nil
}

func (repo *LikeRepository) IsTrackLiked(ctx context.Context,
	userID uuid.UUID,
	trackID uuid.UUID,
) (bool, error) {
	var exists bool
	err := repo.getExecutor(ctx).Get(ctx, &exists, `
		SELECT EXISTS(SELECT 1 FROM liked_tracks WHERE user_id = $1 AND track_id = $2)
	`, userID, trackID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrNotExists
		}
		logger.Logger.Error("failed to check if track is liked",
			slog.String("error", err.Error()),
			slog.String("where", "LikeRepository.IsTrackLiked"))
		return false, ErrSQL
	}

	return exists, nil
}

func (repo *LikeRepository) IsArtistLiked(ctx context.Context,
	userID uuid.UUID,
	artistID uuid.UUID,
) (bool, error) {
	var exists bool
	err := repo.getExecutor(ctx).Get(ctx, &exists, `
		SELECT EXISTS(SELECT 1 FROM liked_artists WHERE user_id = $1 AND artist_id = $2)
	`, userID, artistID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, ErrNotExists
		}
		logger.Logger.Error("failed to check if artist is liked",
			slog.String("error", err.Error()),
			slog.String("where", "LikeRepository.IsArtistLiked"))
		return false, ErrSQL
	}

	return exists, nil
}
