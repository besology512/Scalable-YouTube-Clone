package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Kafka KafkaConfig
	Minio MinioConfig
}

type KafkaConfig struct {
	Brokers string
	GroupID string
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Error reading config file, %s", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return &config, nil
}
