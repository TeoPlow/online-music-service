package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
)

type AlbumRepository struct {
	repository
}

func NewAlbumRepo(db Executor) *AlbumRepository {
	return &AlbumRepository{repository{db}}
}

func (repo *AlbumRepository) Add(ctx context.Context, album models.Album) error {
	_, err := repo.getExecutor(ctx).Exec(ctx, `
		INSERT INTO albums
		VALUES ($1, $2, $3, $4)
	`, album.ID, album.Title, album.ArtistID, album.ReleaseDate)
	if err != nil {
		logger.Logger.Error("failed to insert album",
			slog.String("error", err.Error()),
			slog.String("where", "AlbumRepository.Add"))
		return ErrSQL
	}
	return nil
}

func (repo *AlbumRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := repo.getExecutor(ctx).Exec(ctx, `
		DELETE FROM albums
		WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to delete album",
			slog.String("error", err.Error()),
			slog.String("where", "AlbumRepository.Delete"))
		return ErrSQL
	}
	return nil
}

func (repo *AlbumRepository) Update(ctx context.Context, album models.Album) error {
	_, err := repo.getExecutor(ctx).Exec(ctx, `
		UPDATE albums
		SET
			title = $2,
			artist_id = $3,
			release_date = $4,
		WHERE id = $1
	`, album.ID, album.Title, album.ArtistID, album.ReleaseDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotExists
		}
		logger.Logger.Error("failed to update album",
			slog.String("error", err.Error()),
			slog.String("where", "AlbumRepository.Update"))
		return ErrSQL
	}
	return nil
}

func (repo *AlbumRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Album, error) {
	var res models.Album
	err := repo.getExecutor(ctx).Get(ctx, &res, `
		SELECT *
		FROM albums
		WHERE id = $1
	`, id)
	if err != nil {
		if pgxscan.NotFound(err) {
			return models.Album{}, ErrNotExists
		}
		logger.Logger.Error("failed to get album",
			slog.String("error", err.Error()),
			slog.String("where", "AlbumRepository.GetByID"))
		return models.Album{}, fmt.Errorf("AlbumRepository.GetByID: %w", ErrSQL)
	}
	return res, nil
}

func (repo *AlbumRepository) FindCopy(ctx context.Context, a models.Album) (bool, error) {
	var count int
	if err := repo.getExecutor(ctx).Get(ctx, &count, `
			SELECT COUNT(*)
			FROM albums
			WHERE title = $1 AND artist_id = $2
		`, a.Title, a.ArtistID); err != nil {
		logger.Logger.Error("failed to find copy of album",
			slog.String("error", err.Error()),
			slog.String("where", "AlbumRepository.FindCopy"))
		return false, ErrSQL
	}
	return count > 0, nil
}

func buildAlbumsFilterQuery(offset, limit int, searchQuery *string,
	albumID *uuid.UUID,
) (string, []any) {
	var (
		sb   strings.Builder
		args []any
	)
	argIndex := 1
	sb.WriteString("SELECT * FROM albums WHERE 1=1")

	if searchQuery != nil {
		sb.WriteString(fmt.Sprintf(" AND title LIKE $%d", argIndex))
		args = append(args, "%"+*searchQuery+"%")
		argIndex++
	}

	if albumID != nil {
		sb.WriteString(fmt.Sprintf(" AND artist_id = $%d", argIndex))
		args = append(args, *albumID)
		argIndex++
	}

	sb.WriteString(fmt.Sprintf(" LIMIT $%d", argIndex))
	args = append(args, limit)
	argIndex++
	sb.WriteString(fmt.Sprintf(" OFFSET $%d", argIndex))
	args = append(args, offset)

	return sb.String(), args
}

func (repo *AlbumRepository) GetByFilter(ctx context.Context,
	offset, limit int, searchQuery *string, albumID *uuid.UUID,
) ([]models.Album, error) {
	var res []models.Album
	query, args := buildAlbumsFilterQuery(offset, limit, searchQuery, albumID)
	if err := repo.getExecutor(ctx).Select(ctx, &res, query, args...); err != nil {
		logger.Logger.Error("failed to list albums",
			slog.String("error", err.Error()),
			slog.String("where", "AlbumRepository.GetByID"))
		return nil, ErrSQL
	}
	return res, nil
}
