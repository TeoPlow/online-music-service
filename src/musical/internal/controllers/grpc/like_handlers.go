package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

func (s *Server) LikeTrack(
	ctx context.Context,
	req *musicalpb.LikeTrackRequest,
) (*emptypb.Empty, error) {
	trackID, err := uuid.Parse(req.GetTrackId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = s.likeService.LikeTrack(ctx, userID, trackID)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) UnlikeTrack(
	ctx context.Context,
	req *musicalpb.UnlikeTrackRequest,
) (*emptypb.Empty, error) {
	trackID, err := uuid.Parse(req.GetTrackId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = s.likeService.UnlikeTrack(ctx, userID, trackID)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) LikeArtist(
	ctx context.Context,
	req *musicalpb.LikeArtistRequest,
) (*emptypb.Empty, error) {
	artistID, err := uuid.Parse(req.GetArtistId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = s.likeService.LikeArtist(ctx, userID, artistID)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) UnlikeArtist(
	ctx context.Context,
	req *musicalpb.UnlikeArtistRequest,
) (*emptypb.Empty, error) {
	artistID, err := uuid.Parse(req.GetArtistId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err = s.likeService.UnlikeArtist(ctx, userID, artistID)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetLikedTracks(
	ctx context.Context,
	req *musicalpb.GetLikedTracksRequest,
) (*musicalpb.GetLikedTracksResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	dto, err := protoToGetLikedTracksRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	likedTracks, err := s.likeService.GetLikedTracks(ctx, userID, dto)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	resp := protoFromGetLikedTracksResponse(likedTracks)

	return resp, nil
}

func (s *Server) GetLikedArtists(
	ctx context.Context,
	req *musicalpb.GetLikedArtistsRequest,
) (*musicalpb.GetLikedArtistsResponse, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	dto, err := protoToGetLikedArtistsRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	likedArtists, err := s.likeService.GetLikedArtists(ctx, userID, dto)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	resp := protoFromGetLikedArtistsResponse(likedArtists)
	return resp, nil
}
