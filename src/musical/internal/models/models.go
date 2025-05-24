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

type Album struct {
	ID          uuid.UUID `json:"id"           db:"id"`
	Title       string    `json:"title"        db:"title"`
	ArtistID    uuid.UUID `json:"artist_id"    db:"artist_id"`
	ReleaseDate time.Time `json:"release_date" db:"release_date"`
}

type Artist struct {
	ID          uuid.UUID `json:"id"          db:"id"`
	Name        string    `json:"name"        db:"name"`
	Author      string    `json:"author"      db:"author"`
	Producer    string    `json:"producer"    db:"producer"`
	Country     string    `json:"country"     db:"country"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"  db:"updated_at"`
}
