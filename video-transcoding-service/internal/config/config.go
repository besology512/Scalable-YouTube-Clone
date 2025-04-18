package config

type Config struct {
	MinioEndpoint              string
	MinioAccessKey             string
	MinioSecretKey             string
	MinioRawVideosBucket       string
	MinioProcessedVideosBucket string
	MinioThumbnailsBucket      string
	MinioTmpGOPBucket          string
	MinioTmpGOPEncodedBucket   string
	KafkaBrokers               []string
	KafkaTopic                 string
	KafkaTranscodeTopic        string
	KafkaChunkTopic            string
	KafkaEncodeTopic           string
	KafkaMergeTopic            string
	KafkaUploadTranscodedTopic string
	KafkaTranscodeGroup        string
	KafkaChunkGroup            string
	KafkaEncodeGroup           string
	KafkaMergeGroup            string
	KafkaUploadTranscodedGroup string
	ServerPort                 string
	MongoURI                   string
}

func Load() *Config {
	return &Config{
		MinioEndpoint:              "localhost:9000",
		MinioAccessKey:             "minioadmin",
		MinioSecretKey:             "minioadmin",
		MinioRawVideosBucket:       "raw-videos",
		MinioProcessedVideosBucket: "processed-videos",
		MinioThumbnailsBucket:      "thumbnails",
		MinioTmpGOPBucket:          "tmp-GOP",
		MinioTmpGOPEncodedBucket:   "tmp-GOP-encoded",
		KafkaBrokers:               []string{"localhost:9092"},
		KafkaTopic:                 "video.uploaded",
		KafkaTranscodeTopic:        "video.transcode",
		KafkaChunkTopic:            "video.chunk",
		KafkaEncodeTopic:           "video.encode",
		KafkaMergeTopic:            "video.merge",
		KafkaUploadTranscodedTopic: "video.upload_transcoded",
		KafkaTranscodeGroup:        "transcode-group",
		KafkaChunkGroup:            "chunk-group",
		KafkaEncodeGroup:           "encode-group",
		KafkaMergeGroup:            "merge-group",
		KafkaUploadTranscodedGroup: "upload-transcoded-group",
		ServerPort:                 "8080",
		MongoURI:                   "mongodb://localhost:27017",
	}
}
