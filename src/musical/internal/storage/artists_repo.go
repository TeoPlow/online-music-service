package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
)

type ArtistRepository struct {
	repository
}

func NewArtistRepo(db Executor) *ArtistRepository {
	return &ArtistRepository{repository{db}}
}

func (repo *ArtistRepository) Add(ctx context.Context, artist models.Artist) error {
	_, err := repo.getExecutor(ctx).Exec(ctx, `
		INSERT INTO artists 
		(id, name, author, producer, country, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO NOTHING
	`, artist.ID, artist.Name, artist.Author, artist.Producer,
		artist.Country, artist.Description, artist.CreatedAt, artist.UpdatedAt)
	if err != nil {
		logger.Logger.Error("failed to insert artist",
			slog.String("error", err.Error()),
			slog.String("where", "ArtistRepository.Add"))
		return ErrSQL
	}
	return nil
}

func (repo *ArtistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := repo.getExecutor(ctx).Exec(ctx, `
		DELETE FROM artists
		WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to delete artist",
			slog.String("error", err.Error()),
			slog.String("where", "ArtistRepository.Delete"))
		return ErrSQL
	}
	return nil
}

func (repo *ArtistRepository) Update(ctx context.Context, artist models.Artist) error {
	_, err := repo.getExecutor(ctx).Exec(ctx, `
		UPDATE artists
		SET
			name = $2,
			author = $3,
			producer = $4,
			country = $5,
			description = $6,
			created_at = $7,
			updated_at = $8
		WHERE id = $1
	`, artist.ID, artist.Name, artist.Author, artist.Producer,
		artist.Country, artist.Description, artist.CreatedAt, artist.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to update artist",
			slog.String("error", err.Error()),
			slog.String("where", "ArtistRepository.Update"))
		return ErrSQL
	}
	return nil
}

func (repo *ArtistRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Artist, error) {
	var res models.Artist
	err := repo.getExecutor(ctx).Get(ctx, &res, `
		SELECT *
		FROM artists
		WHERE id = $1
	`, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return models.Artist{}, ErrNotExists
		}
		logger.Logger.Error("failed to get artist",
			slog.String("error", err.Error()),
			slog.String("where", "ArtistRepository.GetByID"))
		return models.Artist{}, ErrSQL
	}
	return res, nil
}

func (repo *ArtistRepository) FindCopy(ctx context.Context, a models.Artist) (bool, error) {
	var count int
	if err := repo.getExecutor(ctx).Get(ctx, &count, `
			SELECT COUNT(*)
			FROM artists
			WHERE name=$1 AND author=$2 AND producer=$3 AND country=$4
		`, a.Name, a.Author, a.Producer, a.Country); err != nil {
		return false, ErrSQL
	}
	return count > 0, nil
}
