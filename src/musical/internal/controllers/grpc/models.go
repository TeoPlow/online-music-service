package grpc

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/TeoPlow/online-music-service/src/musical/internal/models"
	"github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

func protoToTrack(proto *musicalpb.Track) (models.Track, error) {
	id, err := uuid.Parse(proto.GetId())
	if err != nil {
		return models.Track{}, fmt.Errorf("uuid.Parse: %w", err)
	}
	albumID, err := uuid.Parse(proto.GetAlbumId())
	if err != nil {
		return models.Track{}, fmt.Errorf("uuid.Parse: %w", err)
	}
	return models.Track{
		ID:        id,
		Title:     proto.GetTitle(),
		AlbumID:   albumID,
		Genre:     proto.GetGenre(),
		Duration:  uint(proto.GetDuration()),
		CreatedAt: proto.GetCreatedAt().AsTime(),
		UpdatedAt: proto.GetUpdatedAt().AsTime(),
	}, nil
}

func protoFromTrack(track models.Track) *musicalpb.Track {
	return &musicalpb.Track{
		Id:        track.ID.String(),
		Title:     track.Title,
		AlbumId:   track.AlbumID.String(),
		Genre:     track.Genre,
		Duration:  uint32(track.Duration),
		CreatedAt: timestamppb.New(track.CreatedAt),
		UpdatedAt: timestamppb.New(track.UpdatedAt),
	}
}
