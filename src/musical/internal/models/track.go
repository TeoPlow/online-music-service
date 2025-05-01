// Package models определяет основные структуры данных, используемые в бизнес-логике.
package models

import (
	"time"

	"github.com/google/uuid"
)

type Track struct {
	ID        uuid.UUID `json:"id"         db:"id"`
	Title     string    `json:"title"      db:"title"`
	AlbumID   uuid.UUID `json:"album_id"   db:"album_id"`
	Genre     string    `json:"genre"      db:"genre"`
	Duration  uint      `json:"duration"   db:"duration"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
