package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/TeoPlow/online-music-service/src/musical/internal/storage"
	"github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

func (s *Server) AddMusic(ctx context.Context, req *musicalpb.Track) (*emptypb.Empty, error) {
	track, err := protoToTrack(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = s.service.AddMusic(ctx, track); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetMusic(ctx context.Context, req *musicalpb.IDRequest) (*musicalpb.Track, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}

	track, err := s.service.GetMusic(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return protoFromTrack(track), nil
}
