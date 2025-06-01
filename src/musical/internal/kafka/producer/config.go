package producer

import (
	"github.com/IBM/sarama"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
)

func NewSaramaProducerConfig() (*sarama.Config, error) {
	cfg := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(config.Config.Kafka.Version)
	if err != nil {
		return nil, err
	}
	cfg.Version = version

	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Idempotent = true
	cfg.Producer.Retry.Max = config.Config.Kafka.Retries
	cfg.Net.MaxOpenRequests = 1

	return cfg, nil
}
