// Package consumer содержит реализацию Kafka consumer-а для обработки сообщений о событиях, связанных с артистами.
package consumer

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/IBM/sarama"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
)

type KafkaMessageHandler interface {
	HandleMessage(ctx context.Context, topic string, message []byte) error
	GetTopics() []string
}

type ArtistConsumer struct {
	group   sarama.ConsumerGroup
	handler KafkaMessageHandler
	wg      sync.WaitGroup
}

func NewArtistConsumer(
	handler KafkaMessageHandler,
	wg *sync.WaitGroup,
) (*ArtistConsumer, error) {
	saramaCfg, err := NewSaramaConfig()
	if err != nil {
		logger.Logger.Error("failed to create sarama config",
			slog.String("error", err.Error()),
			slog.String("where", "NewArtistConsumer"))
		return nil, err
	}
	group, err := sarama.NewConsumerGroup(
		config.Config.Kafka.Brokers,
		config.Config.Kafka.ConsumerGroup,
		saramaCfg)
	if err != nil {
		logger.Logger.Error("failed to create sarama consumer group",
			slog.String("error", err.Error()),
			slog.String("where", "NewArtistConsumer"))
		return nil, err
	}
	return &ArtistConsumer{
		group:   group,
		handler: handler,
		wg:      sync.WaitGroup{},
	}, nil
}

func (c *ArtistConsumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	logger.Logger.Info("artist consumer claim",
		slog.String("where", "ArtistConsumer.ConsumeClaim"))

	for msg := range claim.Messages() {
		if err := c.handler.HandleMessage(session.Context(), msg.Topic, msg.Value); err != nil {
			logger.Logger.Error("failed to handle message",
				slog.String("error", err.Error()),
				slog.String("where", "ArtistConsumer.ConsumeClaim"))
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *ArtistConsumer) Start(ctx context.Context) error {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {

			if ctx.Err() != nil {
				return
			}

			if err := c.group.Consume(ctx, c.handler.GetTopics(), c); err != nil {
				logger.Logger.Error("failed to consume messages",
					slog.String("error", err.Error()),
					slog.String("where", "ArtistConsumer.Start"))
				select {
				case <-ctx.Done():
					logger.Logger.Info("context done, stopping artist consumer",
						slog.String("where", "ArtistConsumer.Start"))
					return
				case <-time.After(5 * time.Second):
					logger.Logger.Info("retrying to consume messages",
						slog.String("where", "ArtistConsumer.Start"))
				}
				continue
			}
		}
	}()
	return nil
}

func (c *ArtistConsumer) Wait() {
	c.wg.Wait()
}

func (c *ArtistConsumer) Close() error {
	return c.group.Close()
}

func (c *ArtistConsumer) Setup(sarama.ConsumerGroupSession) error {
	logger.Logger.Info("artist consumer setup",
		slog.String("where", "ArtistConsumer.Setup"))
	return nil
}

func (c *ArtistConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	logger.Logger.Info("artist consumer cleanup",
		slog.String("where", "ArtistConsumer.Cleanup"))
	return nil
}
