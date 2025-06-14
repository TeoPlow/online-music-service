// Package config предоставляет функциональность для загрузки и валидации cfg
package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string        `yaml:"env" env-default:"local"`
	GRPC            GRPCConfig    `yaml:"grpc_server"`
	DBURL           string        `yaml:"db_url"`
	JWTSecret       string        `yaml:"jwt_secret"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	LogLevel        string        `yaml:"log_level" env-default:"info"`
	Redis           RedisConfig   `yaml:"redis"`
	Kafka           KafkaConfig   `yaml:"kafka"`
	Outbox          OutboxConfig  `yaml:"outbox"`
	StaticFilesPath string        `yaml:"static_files_path" env:"STATIC_FILES_PATH" default:"static"`
}

type RedisConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	DialTimeout  time.Duration `yaml:"dial_timeout" env-default:"5s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"3s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"3s"`
}

type KafkaConfig struct {
	Brokers        []string      `yaml:"brokers"`
	TLS            TLSConfig     `yaml:"tls"`
	Topics         KafkaTopics   `yaml:"topics"`
	Retries        int           `yaml:"retries"`
	RetryBackoff   time.Duration `yaml:"retry_backoff" env-default:"1s"`
	FlushFrequency time.Duration `yaml:"flush_frequency" env-default:"1s"`
	Compression    string        `yaml:"compression"`
	Version        string        `yaml:"version"`
	ClientID       string        `yaml:"client_id"`
}

type TLSConfig struct {
	Enable      bool          `yaml:"enable"`
	CACert      string        `yaml:"ca_cert"`
	Cert        string        `yaml:"cert"`
	Key         string        `yaml:"key"`
	DialTimeout time.Duration `yaml:"dial_timeout" env-default:"5s"`
}
type OutboxConfig struct {
	BatchSize    int           `yaml:"batch_size"`
	PollInterval time.Duration `yaml:"poll_interval"`
}
type KafkaTopics struct {
	UserCreatedTopic   string `yaml:"user_created"`
	ArtistCreatedTopic string `yaml:"artist_created"`
}
type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func (c RedisConfig) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		panic("Config path is not set")
	}
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config %s", err.Error())
	}
	return &cfg
}
