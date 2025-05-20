// Package dto содержит Data Transfer Objects для всех сервисов
package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateAlbumRequest struct {
	Title    string
	ArtistID uuid.UUID
}

type UpdateAlbumRequest struct {
	ID          uuid.UUID
	Title       *string
	ArtistID    *uuid.UUID
	ReleaseDate *time.Time
}

type ListAlbumsRequest struct {
	Page        int
	PageSize    int
	ArtistID    *uuid.UUID
	SearchQuery *string
}
