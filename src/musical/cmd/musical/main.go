package main

import (
	"context"
	"flag"
	"log"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	"github.com/TeoPlow/online-music-service/src/musical/internal/controllers/grpc"
	"github.com/TeoPlow/online-music-service/src/musical/internal/controllers/http"
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

	storage := storage.NewMusicRepo(tmanager.GetDatabase())
	if err != nil {
		log.Panic(err)
	}

	service := domain.NewMusicService(storage, tmanager)

	httpServer := http.NewServer(service)
	grpcServer := grpc.NewServer(service)

	logger.Logger.Info("Running...")
	go func() {
		if err = grpcServer.Run(); err != nil {
			log.Panic(err)
		}
	}()

	if err = httpServer.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
