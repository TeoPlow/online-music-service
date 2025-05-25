package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

func (s *Server) CreateAlbum(ctx context.Context,
	req *musicalpb.CreateAlbumRequest,
) (*musicalpb.Album, error) {
	dto, err := protoToCreateAlbumRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.albumService.CreateAlbum(ctx, dto)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	return protoFromAlbum(data), nil
}

func (s *Server) UpdateAlbum(ctx context.Context,
	req *musicalpb.UpdateAlbumRequest,
) (*musicalpb.Album, error) {
	dto, err := protoToUpdateAlbumRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.albumService.UpdateAlbum(ctx, dto)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	return protoFromAlbum(data), nil
}

func (s *Server) DeleteAlbum(ctx context.Context,
	req *musicalpb.IDRequest,
) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = s.albumService.DeleteAlbum(ctx, id); err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetAlbum(ctx context.Context,
	req *musicalpb.IDRequest,
) (*musicalpb.Album, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.albumService.GetAlbum(ctx, id)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	return protoFromAlbum(data), nil
}

func (s *Server) ListAlbums(ctx context.Context,
	req *musicalpb.ListAlbumsRequest,
) (*musicalpb.ListAlbumsResponse, error) {
	dto, err := protoToListAlbumsRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	data, err := s.albumService.ListAlbums(ctx, dto)
	if err != nil {
		return nil, status.Error(errorCode(err), err.Error())
	}
	var resp []*musicalpb.Album
	for _, album := range data {
		resp = append(resp, protoFromAlbum(album))
	}
	return &musicalpb.ListAlbumsResponse{Albums: resp}, nil
}
