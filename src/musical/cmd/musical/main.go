package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	"github.com/TeoPlow/online-music-service/src/musical/internal/controllers/grpc"
	"github.com/TeoPlow/online-music-service/src/musical/internal/db"
	"github.com/TeoPlow/online-music-service/src/musical/internal/domain"
	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
	"github.com/TeoPlow/online-music-service/src/musical/internal/storage"
)

func main() {
	filepath := flag.String("path", "configs/conf.yml", "путь к конфигурационному файлу")
	flag.Parse()

	if err := config.Load(*filepath); err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.InitLogger()

	tmanager, err := db.NewTxManager(ctx)
	if err != nil {
		log.Panic(err)
	}
	defer tmanager.GetDatabase().Close()

	artistRepo := storage.NewArtistRepo(tmanager.GetDatabase())
	artists := domain.NewArtistService(artistRepo, tmanager)

	albumRepo := storage.NewAlbumRepo(tmanager.GetDatabase())
	albums := domain.NewAlbumService(albumRepo, tmanager, artists)

	grpcServer := grpc.NewServer(artists, albums)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Logger.Info("Running...")
		if err := grpcServer.Run(); err != nil {
			log.Panic(err)
		}
	}()

	<-stop
	logger.Logger.Info("Shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := grpcServer.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Error("failed to shutdown", slog.String("error", err.Error()))
	} else {
		logger.Logger.Info("Server stopped cleanly")
	}
}
