// Package producer содержит реализацию Kafka producer-а с использованием библиотеки Sarama.
package producer

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/IBM/sarama"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
)

type KafkaProducer interface {
	Publish(
		ctx context.Context,
		topic string,
		key string,
		payload []byte,
		eventID int64,
	) error
	Close()
}

type SaramaProducer struct {
	asyncProducer sarama.AsyncProducer
	topics        struct {
		ArtistCreated string `yaml:"artist_created"`
		ArtistUpdated string `yaml:"artist_updated"`
		ArtistDeleted string `yaml:"artist_deleted"`
	}
	wg sync.WaitGroup
}

func NewSaramaProducer(
	cfg config.KafkaConfig,
	topics struct {
		ArtistCreated string `yaml:"artist_created"`
		ArtistUpdated string `yaml:"artist_updated"`
		ArtistDeleted string `yaml:"artist_deleted"`
	},
) (*SaramaProducer, error) {
	saramaCfg, err := NewSaramaProducerConfig(cfg)
	if err != nil {
		return nil, err
	}

	asyncProducer, err := sarama.NewAsyncProducer(cfg.Brokers, saramaCfg)
	if err != nil {
		return nil, err
	}

	producer := &SaramaProducer{
		asyncProducer: asyncProducer,
		topics:        topics,
	}

	producer.wg.Add(2)
	go producer.handleErrors()
	go producer.handleSuccesses()

	return producer, nil
}

func (p *SaramaProducer) handleSuccesses() {
	const op = "kafka.SaramaProducer.handleSuccesses"
	defer p.wg.Done()

	for msg := range p.asyncProducer.Successes() {
		logger.Logger.Info("message successfully sent",
			slog.String("topic", msg.Topic),
			slog.Int64("offset", msg.Offset),
			slog.String("where", op))

	}
}

func (p *SaramaProducer) handleErrors() {
	const op = "kafka.SaramaProducer.handleErrors"
	defer p.wg.Done()

	for err := range p.asyncProducer.Errors() {
		logger.Logger.Error("failed to send message",
			slog.String("error", err.Error()),
			slog.String("topic", err.Msg.Topic),
			slog.String("where", op))
	}
}

func (p *SaramaProducer) Publish(
	ctx context.Context,
	topic string,
	key string,
	payload []byte,
	eventID int64,
) error {
	const op = "kafka.SaramaProducer.Publish"

	if topic == "" {
		return errors.New("topic cannot be empty")
	}

	if len(payload) == 0 {
		return errors.New("payload cannot be empty")
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(payload),
		Metadata:  eventID,
		Timestamp: time.Now(),
	}

	select {
	case p.asyncProducer.Input() <- msg:
		logger.Logger.Info("message published to Kafka",
			slog.String("topic", topic),
			slog.String("key", key),
			slog.String("where", op))
		return nil
	case err := <-p.asyncProducer.Errors():
		logger.Logger.Error("failed to publish message",
			slog.String("error", err.Error()),
			slog.String("topic", topic),
			slog.String("where", op))
		return err
	}
}

func (p *SaramaProducer) Close() {
	const op = "kafka.SaramaProducer.Close"
	if err := p.asyncProducer.Close(); err != nil {
		logger.Logger.Error("failed to close producer",
			slog.String("error", err.Error()),
			slog.String("where", op))
	}
	p.wg.Wait()
	logger.Logger.Info("producer closed successfully", slog.String("where", op))
}
