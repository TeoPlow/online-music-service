// Package config предоставляет структуру и позволяет загружать конфигурацию приложения
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	DBConn   string `yaml:"dbconnection"`
	GRPCPort string `yaml:"grpc_port"`
	Log      struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"log"`
	Minio struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"access"`
		SecretKey string `yaml:"secret"`
		UseSSL    bool   `yaml:"ssl"`
	} `yaml:"minio"`
}

var Config config

func Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &Config); err != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	return nil
}
