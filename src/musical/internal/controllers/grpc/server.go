// Package grpc реализует gRPC-обработчики для удалённого взаимодействия с сервисом.
package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	"github.com/TeoPlow/online-music-service/src/musical/internal/controllers"
	"github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

type Server struct {
	musicalpb.UnimplementedMusicalServiceServer
	artistService    controllers.ArtistService
	albumService     controllers.AlbumService
	trackService     controllers.TrackService
	streamingService controllers.StreamingService

	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(
	artistServ controllers.ArtistService,
	albumServ controllers.AlbumService,
	trackServ controllers.TrackService,
	streamServ controllers.StreamingService,
) *Server {
	return &Server{
		artistService:    artistServ,
		albumService:     albumServ,
		trackService:     trackServ,
		streamingService: streamServ,
	}
}

func (s *Server) Run() error {
	s.grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(
		logInterceptor,
	))

	var err error
	s.listener, err = net.Listen("tcp", config.Config.GRPCPort)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	s.Register(s.grpcServer)
	reflection.Register(s.grpcServer)

	if err := s.grpcServer.Serve(s.listener); err != nil {
		return fmt.Errorf("grpc.Server.Serve: %w", err)
	}
	return nil
}

func (s *Server) Register(grpcServer *grpc.Server) {
	musicalpb.RegisterMusicalServiceServer(grpcServer, s)
}

func (s *Server) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
		return ctx.Err()
	case <-stopped:
		return nil
	}
}
