package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
)

type trackForDB struct {
	ID         uuid.UUID      `db:"id"`
	Title      string         `db:"title"`
	AlbumID    uuid.UUID      `db:"album_id"`
	Genre      string         `db:"genre"`
	Duration   time.Duration  `db:"duration"`
	Lyrics     sql.NullString `db:"lyrics"`
	IsExplicit bool           `db:"is_explicit"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
}

func convertToDB(t models.Track) trackForDB {
	var lyrics sql.NullString
	if t.Lyrics != nil {
		lyrics.Valid = true
		lyrics.String = *t.Lyrics
	}
	return trackForDB{
		ID:         t.ID,
		Title:      t.Title,
		AlbumID:    t.AlbumID,
		Genre:      t.Genre,
		Duration:   t.Duration,
		Lyrics:     lyrics,
		IsExplicit: t.IsExplicit,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

func convertToModel(t trackForDB) models.Track {
	var lyrics *string
	if t.Lyrics.Valid {
		lyrics = &t.Lyrics.String
	}
	return models.Track{
		ID:         t.ID,
		Title:      t.Title,
		AlbumID:    t.AlbumID,
		Genre:      t.Genre,
		Duration:   t.Duration,
		Lyrics:     lyrics,
		IsExplicit: t.IsExplicit,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

type TrackRepository struct {
	repository
}

func NewTrackRepo(db Executor) *TrackRepository {
	return &TrackRepository{repository{db}}
}

func (repo *TrackRepository) Add(ctx context.Context, t models.Track) error {
	track := convertToDB(t)
	_, err := repo.getExecutor(ctx).Exec(ctx, `
		INSERT INTO tracks
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, track.ID, track.Title, track.AlbumID, track.Genre,
		track.Duration, track.Lyrics, track.IsExplicit,
		track.CreatedAt, track.UpdatedAt)
	if err != nil {
		logger.Logger.Error("failed to insert track",
			slog.String("error", err.Error()),
			slog.String("where", "TrackRepository.Add"))
		return ErrSQL
	}
	return nil
}

func (repo *TrackRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := repo.getExecutor(ctx).Exec(ctx, `
		DELETE FROM tracks
		WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to delete track",
			slog.String("error", err.Error()),
			slog.String("where", "TrackRepository.Delete"))
		return ErrSQL
	}
	return nil
}

func (repo *TrackRepository) Update(ctx context.Context, t models.Track) error {
	track := convertToDB(t)
	_, err := repo.getExecutor(ctx).Exec(ctx, `
		UPDATE tracks
		SET
			title = $2,
			album_id = $3,
			genre = $4,
			duration = $5,
			lyrics = $6,
			is_explicit = $7,
			updated_at = $8
		WHERE id = $1
	`, track.ID, track.Title, track.AlbumID, track.Genre,
		track.Duration, track.Lyrics, track.IsExplicit,
		track.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to update track",
			slog.String("error", err.Error()),
			slog.String("where", "TrackRepository.Update"))
		return ErrSQL
	}
	return nil
}

func (repo *TrackRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Track, error) {
	var res trackForDB
	err := repo.getExecutor(ctx).Get(ctx, &res, `
		SELECT *
		FROM tracks
		WHERE id = $1
	`, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return models.Track{}, ErrNotExists
		}
		logger.Logger.Error("failed to get track",
			slog.String("error", err.Error()),
			slog.String("where", "TrackRepository.GetByID"))
		return models.Track{}, fmt.Errorf("TrackRepository.GetByID: %w", ErrSQL)
	}
	return convertToModel(res), nil
}

func (repo *TrackRepository) FindCopy(ctx context.Context, a models.Track) (bool, error) {
	var count int
	if err := repo.getExecutor(ctx).Get(ctx, &count, `
			SELECT COUNT(*)
			FROM tracks
			WHERE title = $1 AND album_id = $2
		`, a.Title, a.AlbumID); err != nil {
		logger.Logger.Error("failed to find copy of track",
			slog.String("error", err.Error()),
			slog.String("where", "TrackRepository.FindCopy"))
		return false, ErrSQL
	}
	return count > 0, nil
}

func buildTracksFilterQuery(offset, limit int, genre, searchQuery *string,
	trackID *uuid.UUID,
) (string, []any) {
	var (
		sb   strings.Builder
		args []any
	)
	argIndex := 1
	sb.WriteString("SELECT * FROM tracks WHERE 1=1")

	if genre != nil {
		sb.WriteString(fmt.Sprintf(" AND genre = $%d", argIndex))
		args = append(args, *genre)
		argIndex++
	}

	if searchQuery != nil {
		sb.WriteString(fmt.Sprintf(
			" AND (similarity(title, $%d) > %f OR (lyrics IS NOT NULL AND similarity(lyrics, $%d) > %f))",
			argIndex, threshold, argIndex, threshold))
		args = append(args, *searchQuery)
		argIndex++
	}

	if trackID != nil {
		sb.WriteString(fmt.Sprintf(" AND album_id = $%d", argIndex))
		args = append(args, *trackID)
		argIndex++
	}

	if searchQuery != nil {
		sb.WriteString(fmt.Sprintf(
			" ORDER BY GREATEST(similarity(title, $%d), COALESCE(similarity(lyrics, $%d), 0)) DESC",
			argIndex, argIndex))
		args = append(args, *searchQuery)
		argIndex++
	}

	sb.WriteString(fmt.Sprintf(" LIMIT $%d", argIndex))
	args = append(args, limit)
	argIndex++
	sb.WriteString(fmt.Sprintf(" OFFSET $%d", argIndex))
	args = append(args, offset)

	return sb.String(), args
}

func (repo *TrackRepository) GetByFilter(ctx context.Context,
	offset, limit int, genre, searchQuery *string, trackID *uuid.UUID,
) ([]models.Track, error) {
	var res []models.Track
	query, args := buildTracksFilterQuery(offset, limit, genre, searchQuery, trackID)
	if err := repo.getExecutor(ctx).Select(ctx, &res, query, args...); err != nil {
		logger.Logger.Error("failed to list tracks",
			slog.String("error", err.Error()),
			slog.String("where", "TrackRepository.GetByID"))
		return nil, ErrSQL
	}
	return res, nil
}
