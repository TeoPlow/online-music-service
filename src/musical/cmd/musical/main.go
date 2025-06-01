package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	"github.com/TeoPlow/online-music-service/src/musical/internal/controllers"
	"github.com/TeoPlow/online-music-service/src/musical/internal/controllers/grpc"
	"github.com/TeoPlow/online-music-service/src/musical/internal/db"
	"github.com/TeoPlow/online-music-service/src/musical/internal/domain"
	"github.com/TeoPlow/online-music-service/src/musical/internal/kafka/consumer"
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

	minio, err := storage.NewMinIOClient(ctx, "audio")
	if err != nil {
		log.Panic(err)
	}

	if err != nil {
		logger.Logger.Error("failed to create Kafka producer",
			slog.String("error", err.Error()),
			slog.String("where", "main.NewSaramaProducer"))
		return
	}
	artistRepo := storage.NewArtistRepo(tmanager.GetDatabase())
	artists := domain.NewArtistService(artistRepo, tmanager)

	albumRepo := storage.NewAlbumRepo(tmanager.GetDatabase())
	albums := domain.NewAlbumService(albumRepo, tmanager, artists)

	streaming := domain.NewStreamingService(minio)

	trackRepo := storage.NewTrackRepo(tmanager.GetDatabase())
	tracks := domain.NewTrackService(trackRepo, tmanager, albums, streaming)

	likeRepo := storage.NewLikeRepo(tmanager.GetDatabase())
	likes := domain.NewLikeService(likeRepo, artists, tracks, tmanager)
	wg := &sync.WaitGroup{}
	if err := setupKafkaConsumer(ctx, config.Config.Kafka, *artists, wg); err != nil {
		logger.Logger.Error("failed to setup kafka consumer",
			slog.String("error", err.Error()),
			slog.String("where", "main.setupKafkaConsumer"))
		return
	}

	grpcServer := grpc.NewServer(artists, albums, tracks, streaming, likes)

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

func setupKafkaConsumer(ctx context.Context,
	cfg config.KafkaConfig,
	artistService domain.ArtistService,
	wg *sync.WaitGroup,
) error {
	kafkaHandler := controllers.NewKafkaArtistHandler(artistService, config.Config.Kafka.Topics)

	artistConsumer, err := consumer.NewArtistConsumer(cfg, kafkaHandler, wg)
	if err != nil {
		logger.Logger.Error("failed to create artist consumer",
			slog.String("error", err.Error()),
			slog.String("where", "setupKafkaConsumer"))
		return err
	}
	if err := artistConsumer.Start(ctx); err != nil {
		logger.Logger.Error("failed to start artist consumer",
			slog.String("error", err.Error()),
			slog.String("where", "setupKafkaConsumer"))
		return err
	}
	go func() {
		<-ctx.Done()
		artistConsumer.Wait()
		if err := artistConsumer.Close(); err != nil {
			logger.Logger.Error("failed to close artist consumer",
				slog.String("error", err.Error()),
				slog.String("where", "setupKafkaConsumer"))
		} else {
			logger.Logger.Info("artist consumer closed successfully")
		}
	}()

	return nil
}
