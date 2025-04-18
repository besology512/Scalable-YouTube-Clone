package clients

import (
	"net/http"
	"time"
)

type FunctionClient struct {
	BaseURL string
	HTTP    *http.Client
}

func (a *FunctionClient) GetBaseURL() string {
	return a.BaseURL
}
func NewFunctionClient(baseURL string) *FunctionClient {
	return &FunctionClient{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}
