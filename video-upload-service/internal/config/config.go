package config

type Config struct {
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	KafkaBrokers   string
	KafkaTopic     string
	ServerPort     string
	MongoURI       string
}

func Load() *Config {
	return &Config{
		MinioEndpoint:  "localhost:9000",
		MinioAccessKey: "minioadmin",
		MinioSecretKey: "minioadmin",
		MinioBucket:    "raw-videos",
		KafkaBrokers:   "localhost:9092",
		KafkaTopic:     "video.uploaded",
		ServerPort:     "8080",
		MongoURI:       "mongodb://localhost:27017",
	}
}
