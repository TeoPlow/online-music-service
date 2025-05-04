package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
)

type DB interface {
	Select(ctx context.Context, dst any, query string, args ...any) error
	Get(ctx context.Context, dst any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type MusicRepository struct {
	db DB
}

func NewMusicRepo(db DB) *MusicRepository {
	return &MusicRepository{db: db}
}

func (s *MusicRepository) Add(ctx context.Context, track models.Track) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO tracks
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`, track.ID, track.Title, track.AlbumID,
		track.Genre, track.Duration,
		track.CreatedAt, track.UpdatedAt)
	if err != nil {
		logger.Logger.Error(err.Error(), slog.String("where", "MusicRepository.Add"))
		return fmt.Errorf("MusicRepository.Add: %w", ErrSQL)
	}
	return nil
}

func (s *MusicRepository) Delete(ctx context.Context, id uuid.UUID) error {
	logger.Logger.Error("Unimplemented", slog.String("where", "MusicRepository.Delete"))
	return nil
}

func (s *MusicRepository) Update(ctx context.Context, t models.Track) error {
	logger.Logger.Error("Unimplemented", slog.String("where", "MusicRepository.Update"))
	return nil
}

func (s *MusicRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Track, error) {
	var res models.Track
	err := s.db.Get(ctx, &res, `
		SELECT *
		FROM tracks
		WHERE id = $1
	`, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return models.Track{},
				fmt.Errorf(
					"MusicRepository.GetByID: track with id=%s %w",
					id.String(), ErrNotExists)
		}
		logger.Logger.Error(err.Error(), slog.String("where", "MusicRepository.GetByID"))
		return models.Track{}, fmt.Errorf("MusicRepository.GetByID: %w", ErrSQL)
	}
	return res, nil
}

func (s *MusicRepository) GetAll(ctx context.Context) ([]models.Track, error) {
	logger.Logger.Error("Unimplemented", slog.String("where", "MusicRepository.GetAll"))
	return []models.Track{}, nil
}
