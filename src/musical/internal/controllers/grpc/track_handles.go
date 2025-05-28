package grpc

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

func handleCreateStream(stream musicalpb.MusicalService_CreateTrackServer,
) (*musicalpb.TrackInfo, *bytes.Buffer, error) {
	var (
		trackInfo *musicalpb.TrackInfo
		buffer    bytes.Buffer
	)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, &bytes.Buffer{},
				status.Error(codes.Internal, "failed to receive stream: "+err.Error())
		}

		switch x := req.Data.(type) {
		case *musicalpb.CreateTrackRequest_Metadata:
			if trackInfo != nil {
				return nil, &bytes.Buffer{},
					status.Error(codes.InvalidArgument, "metadata already received")
			}
			trackInfo = x.Metadata

		case *musicalpb.CreateTrackRequest_Chunk:
			if trackInfo == nil {
				return nil, &bytes.Buffer{},
					status.Error(codes.InvalidArgument, "metadata must be sent before file chunks")
			}
			_, err := buffer.Write(x.Chunk)
			if err != nil {
				return nil, &bytes.Buffer{},
					status.Error(codes.Internal, "failed to write chunk: "+err.Error())
			}

		default:
			return nil, &bytes.Buffer{},
				status.Error(codes.InvalidArgument, "unknown stream data type")
		}
	}

	if trackInfo == nil {
		return nil, &bytes.Buffer{},
			status.Error(codes.InvalidArgument, "metadata was not sent")
	}
	return trackInfo, &buffer, nil
}

func (s *Server) CreateTrack(stream musicalpb.MusicalService_CreateTrackServer) error {
	trackInfo, buffer, err := handleCreateStream(stream)
	if err != nil {
		return err
	}

	dto, err := protoToCreateTrackRequest(trackInfo)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	track, err := s.trackService.CreateTrack(stream.Context(), dto, buffer)
	if err != nil {
		return status.Error(errorCode(err), "failed to store track file: "+err.Error())
	}

	return stream.SendAndClose(protoFromTrack(track))
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

func (s *Server) DownloadTrack(
	req *musicalpb.IDRequest,
	stream musicalpb.MusicalService_DownloadTrackServer,
) error {
	ctx := stream.Context()

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	object, err := s.streamingService.DownloadTrack(ctx, id)
	if err != nil {
		return status.Error(errorCode(err), err.Error())
	}
	defer func() {
		if err := object.Close(); err != nil {
			logger.Logger.Error("failed close file",
				slog.String("error", err.Error()))
		}
	}()
	const chunkSize = 1024 * 32 // 32 KB
	buf := make([]byte, chunkSize)
	for {
		n, err := object.Read(buf)
		if err != nil && err != io.EOF {
			return status.Error(codes.Internal, "error reading object: "+err.Error())
		}
		if n == 0 {
			break
		}

		if err := stream.Send(&musicalpb.DownloadResponse{
			Chunk: buf[:n],
		}); err != nil {
			return status.Error(codes.Internal, "error sending chunk: "+err.Error())
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}
