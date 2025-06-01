package consumer

import (
	"github.com/IBM/sarama"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
)

func NewSaramaConfig() (*sarama.Config, error) {
	cfg := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(config.Config.Kafka.Version)
	if err != nil {
		return nil, err
	}
	cfg.Version = version

	cfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Offsets.AutoCommit.Enable = true
	cfg.Consumer.Offsets.AutoCommit.Interval = config.Config.Kafka.AutoCommitInterval

	return cfg, nil
}
