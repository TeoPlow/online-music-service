// Package grpcapp инициализирует и запускает gRPC-сервер.
package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/TeoPlow/online-music-service/src/auth/internal/grpc/interceptors"
	"github.com/TeoPlow/online-music-service/src/auth/internal/logger"
	"github.com/TeoPlow/online-music-service/src/auth/internal/pkg/jwt"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
}

func New(port int, jwtCfg jwt.Config) *App {
	authInterceptor := interceptors.NewAuthInterceptor(jwtCfg)
	loggingInterceptor := interceptors.NewLoggingInterceptor()

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor,
			authInterceptor,
		),
	)

	return &App{
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Log.With(slog.String("op", op)).Info(op, "grpc server is running ", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	logger.Log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}

func (a *App) GetServer() *grpc.Server {
	return a.gRPCServer
}
