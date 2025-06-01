package producer

import (
	"github.com/IBM/sarama"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
)

func NewSaramaProducerConfig(
	kcfg config.KafkaConfig,
) (*sarama.Config, error) {
	cfg := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(kcfg.Version)
	if err != nil {
		return nil, err
	}
	cfg.Version = version

	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Idempotent = true
	cfg.Producer.Retry.Max = kcfg.Retries

	return cfg, nil
}
