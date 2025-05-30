// Package controllers содержит основные контроллеры приложения (http и grpc)
package controllers

import (
	"bytes"
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models/dto"
)

type StreamingService interface {
	DownloadTrack(context.Context, uuid.UUID) (io.ReadCloser, error)
}

type TrackService interface {
	CreateTrack(context.Context, dto.CreateTrackRequest, *bytes.Buffer) (models.Track, error)
	UpdateTrack(context.Context, dto.UpdateTrackRequest) (models.Track, error)
	DeleteTrack(ctx context.Context, id uuid.UUID) error
	GetTrack(ctx context.Context, id uuid.UUID) (models.Track, error)
	ListTracks(context.Context, dto.ListTracksRequest) ([]models.Track, error)
}

type AlbumService interface {
	CreateAlbum(context.Context, dto.CreateAlbumRequest) (models.Album, error)
	UpdateAlbum(context.Context, dto.UpdateAlbumRequest) (models.Album, error)
	DeleteAlbum(ctx context.Context, id uuid.UUID) error
	GetAlbum(ctx context.Context, id uuid.UUID) (models.Album, error)
	ListAlbums(context.Context, dto.ListAlbumsRequest) ([]models.Album, error)
}

type ArtistService interface {
	GetArtist(ctx context.Context, id uuid.UUID) (models.Artist, error)
}
