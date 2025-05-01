// Package grpc реализует gRPC-обработчики для удалённого взаимодействия с сервисом.
package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	"github.com/TeoPlow/online-music-service/src/musical/internal/controllers"
	"github.com/TeoPlow/online-music-service/src/musical/pkg/musicalpb"
)

type Server struct {
	musicalpb.UnimplementedMusicServiceServer
	service controllers.MusicService
}

func NewServer(s controllers.MusicService) *Server {
	return &Server{service: s}
}

func (s *Server) Run() error {
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logInterceptor,
	))

	listener, err := net.Listen("tcp", config.Config.GRPCPort)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	musicalpb.RegisterMusicServiceServer(server, s)

	if err = server.Serve(listener); err != nil {
		return fmt.Errorf("grpc.Server.Serve: %w", err)
	}
	return nil
}
