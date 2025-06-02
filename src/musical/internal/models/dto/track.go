package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTrackRequest struct {
	Title      string
	AlbumID    uuid.UUID
	Genre      string
	Duration   time.Duration
	Lyrics     *string
	IsExplicit bool
}

type UpdateTrackRequest struct {
	ID         uuid.UUID
	Title      *string
	AlbumID    *uuid.UUID
	Genre      *string
	Duration   *time.Duration
	Lyrics     *string
	IsExplicit *bool
}

type ListTracksRequest struct {
	Page        int
	PageSize    int
	AlbumID     *uuid.UUID
	Genre       *string
	SearchQuery *string
}
