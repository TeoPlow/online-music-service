package grpc

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/TeoPlow/online-music-service/src/musical/internal/domain"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	"github.com/TeoPlow/online-music-service/src/musical/internal/models/dto"
	pb "github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

func errorCode(err error) codes.Code {
	switch {
	case errors.Is(err, domain.ErrAlreadyExists):
		return codes.AlreadyExists
	case errors.Is(err, domain.ErrNotFound):
		return codes.NotFound
	case errors.Is(err, domain.ErrInternal):
		return codes.Internal
	default:
		return codes.Internal
	}
}

func protoFromAlbum(a models.Album) *pb.Album {
	return &pb.Album{
		Id:          a.ID.String(),
		Title:       a.Title,
		ArtistId:    a.ArtistID.String(),
		ReleaseDate: timestamppb.New(a.ReleaseDate),
	}
}

func protoToCreateAlbumRequest(p *pb.CreateAlbumRequest) (dto.CreateAlbumRequest, error) {
	artistID, err := uuid.Parse(p.GetArtistId())
	if err != nil {
		return dto.CreateAlbumRequest{}, err
	}

	return dto.CreateAlbumRequest{
		Title:    p.GetTitle(),
		ArtistID: artistID,
	}, nil
}

func protoToUpdateAlbumRequest(p *pb.UpdateAlbumRequest) (dto.UpdateAlbumRequest, error) {
	id, err := uuid.Parse(p.GetId())
	if err != nil {
		return dto.UpdateAlbumRequest{}, err
	}

	req := dto.UpdateAlbumRequest{
		ID: id,
	}

	if p.Title != nil {
		req.Title = p.Title
	}

	if p.ArtistId != nil {
		artistID, err := uuid.Parse(*p.ArtistId)
		if err != nil {
			return dto.UpdateAlbumRequest{}, err
		}
		req.ArtistID = &artistID
	}

	if p.ReleaseDate != nil {
		t := p.GetReleaseDate().AsTime()
		req.ReleaseDate = &t
	}

	return req, nil
}

func protoToListAlbumsRequest(p *pb.ListAlbumsRequest) (dto.ListAlbumsRequest, error) {
	req := dto.ListAlbumsRequest{
		Page:     int(p.GetPage()),
		PageSize: int(p.GetPageSize()),
	}

	if p.ArtistId != nil {
		artistID, err := uuid.Parse(*p.ArtistId)
		if err != nil {
			return dto.ListAlbumsRequest{}, err
		}
		req.ArtistID = &artistID
	}

	if p.SearchQuery != nil {
		req.SearchQuery = p.SearchQuery
	}

	return req, nil
}

func protoFromArtist(a models.Artist) *pb.Artist {
	return &pb.Artist{
		Id:          a.ID.String(),
		Name:        a.Name,
		Author:      a.Author,
		Producer:    a.Producer,
		Country:     a.Country,
		Description: a.Description,
		CreatedAt:   timestamppb.New(a.CreatedAt),
		UpdatedAt:   timestamppb.New(a.UpdatedAt),
	}
}

func protoFromTrack(t models.Track) *pb.Track {
	return &pb.Track{
		Id:         t.ID.String(),
		Title:      t.Title,
		AlbumId:    t.AlbumID.String(),
		Genre:      t.Genre,
		Duration:   int32(t.Duration.Seconds()),
		Lyrics:     t.Lyrics,
		IsExplicit: t.IsExplicit,
		CreatedAt:  timestamppb.New(t.CreatedAt),
		UpdatedAt:  timestamppb.New(t.UpdatedAt),
	}
}

func protoToCreateTrackRequest(p *pb.CreateTrackRequest) (dto.CreateTrackRequest, error) {
	albumID, err := uuid.Parse(p.GetAlbumId())
	if err != nil {
		return dto.CreateTrackRequest{}, err
	}
	return dto.CreateTrackRequest{
		Title:      p.Title,
		AlbumID:    albumID,
		Genre:      p.Genre,
		Duration:   time.Duration(p.GetDuration()) * time.Second,
		Lyrics:     p.Lyrics,
		IsExplicit: p.IsExplicit,
	}, nil
}

func protoToUpdateTrackRequest(p *pb.UpdateTrackRequest) (dto.UpdateTrackRequest, error) {
	albumID, err := uuid.Parse(p.GetAlbumId())
	if err != nil {
		return dto.UpdateTrackRequest{}, err
	}
	id, err := uuid.Parse(p.GetId())
	if err != nil {
		return dto.UpdateTrackRequest{}, err
	}
	duration := time.Duration(p.GetDuration()) * time.Second
	return dto.UpdateTrackRequest{
		ID:         id,
		Title:      p.Title,
		AlbumID:    &albumID,
		Genre:      p.Genre,
		Duration:   &duration,
		Lyrics:     p.Lyrics,
		IsExplicit: p.IsExplicit,
	}, nil
}

func protoToListTracksRequest(p *pb.ListTracksRequest) (dto.ListTracksRequest, error) {
	req := dto.ListTracksRequest{
		Page:        int(p.GetPage()),
		PageSize:    int(p.GetPageSize()),
		Genre:       p.Genre,
		SearchQuery: p.SearchQuery,
	}

	if p.AlbumId != nil {
		albumID, err := uuid.Parse(*p.AlbumId)
		if err != nil {
			return dto.ListTracksRequest{}, err
		}
		req.AlbumID = &albumID
	}

	return req, nil
}
