package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

func (s *Server) CreateTrack(ctx context.Context,
	req *musicalpb.CreateTrackRequest,
) (*musicalpb.Track, error) {
	dto, err := protoToCreateTrackRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.trackService.CreateTrack(ctx, dto)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	return protoFromTrack(data), nil
}

func (s *Server) UpdateTrack(ctx context.Context,
	req *musicalpb.UpdateTrackRequest,
) (*musicalpb.Track, error) {
	dto, err := protoToUpdateTrackRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.trackService.UpdateTrack(ctx, dto)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	return protoFromTrack(data), nil
}

func (s *Server) DeleteTrack(ctx context.Context,
	req *musicalpb.IDRequest,
) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = s.trackService.DeleteTrack(ctx, id); err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetTrack(ctx context.Context,
	req *musicalpb.IDRequest,
) (*musicalpb.Track, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.trackService.GetTrack(ctx, id)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	return protoFromTrack(data), nil
}

func (s *Server) ListTracks(ctx context.Context,
	req *musicalpb.ListTracksRequest,
) (*musicalpb.ListTracksResponse, error) {
	dto, err := protoToListTracksRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.trackService.ListTracks(ctx, dto)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	var resp []*musicalpb.Track
	for _, track := range data {
		resp = append(resp, protoFromTrack(track))
	}
	return &musicalpb.ListTracksResponse{Tracks: resp}, nil
}
