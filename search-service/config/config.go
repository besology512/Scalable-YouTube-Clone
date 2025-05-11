package config

import (
	"os"
)

type Config struct {
	ElasticsearchURL string
}

func LoadConfig() Config {
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://localhost:9200"
	}
	return Config{ElasticsearchURL: url}
}
