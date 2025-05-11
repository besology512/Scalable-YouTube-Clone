package main

import (
	"os"
	"search-service/config"
	"search-service/elastic"
	"search-service/router"
	"search-service/sync"
)

func main() {
	cfg := config.LoadConfig()
	elastic.InitClient(cfg.ElasticsearchURL)

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	sync.SyncFromMongo(mongoURI)

	r := router.SetupRouter()
	r.Run(":8080") // Start server
}
