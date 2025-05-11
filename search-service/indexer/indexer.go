package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"search-service/elastic"
	"search-service/model"
)

const IndexName = "videos"

func IndexVideo(video model.Video) error {
	body, _ := json.Marshal(video)
	req := bytes.NewReader(body)

	_, err := elastic.ES.Index(
		IndexName,
		req,
		elastic.ES.Index.WithDocumentID(video.VideoID),
		elastic.ES.Index.WithContext(context.Background()),
	)
	return err
}
