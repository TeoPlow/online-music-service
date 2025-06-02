package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

func (s *Server) GetArtist(ctx context.Context,
	req *musicalpb.IDRequest,
) (*musicalpb.Artist, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.artistService.GetArtist(ctx, id)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	return protoFromArtist(data), nil
}
