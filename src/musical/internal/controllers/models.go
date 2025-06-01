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
	DownloadTrack(context.Context, uuid.UUID, uuid.UUID) (io.ReadCloser, error)
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
	CreateArtist(ctx context.Context, artist models.Artist) error
	UpdateArtist(ctx context.Context, artist models.Artist) error
	DeleteArtist(ctx context.Context, id uuid.UUID) error
}

type LikeService interface {
	LikeTrack(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error
	UnlikeTrack(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error
	LikeArtist(ctx context.Context, userID uuid.UUID, artistID uuid.UUID) error
	UnlikeArtist(ctx context.Context, userID uuid.UUID, artistID uuid.UUID) error
	GetLikedArtists(ctx context.Context, userID uuid.UUID, req dto.GetLikedArtistsRequest) ([]models.Artist, error)
	GetLikedTracks(ctx context.Context, userID uuid.UUID, req dto.GetLikedTracksRequest) ([]models.Track, error)
}

type KafkaMessageHandler interface {
	HandleMessage(ctx context.Context, topic string, message []byte) error
}
