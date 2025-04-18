package handlers

import (
	"io"
	"net/http"

	"github.com/besology512/api-gateway/internal/gateway/clients"
	"github.com/gin-gonic/gin"
)

type FunctionHandler struct {
	client *clients.FunctionClient // holds base URL and HTTP client
}

func NewFunctionHandler(c *clients.FunctionClient) *FunctionHandler {
	return &FunctionHandler{client: c}
}

func (request *FunctionHandler) Proxy(GinContext *gin.Context) {

	target := request.client.GetBaseURL() + GinContext.Request.URL.Path

	if q := GinContext.Request.URL.RawQuery; q != "" {
		target += "?" + q
	}

	req, err := http.NewRequestWithContext(
		GinContext.Request.Context(),
		GinContext.Request.Method,
		target,
		GinContext.Request.Body,
	)

	if err != nil {
		GinContext.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	for k, v := range GinContext.Request.Header {
		req.Header[k] = v
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		GinContext.JSON(http.StatusServiceUnavailable, gin.H{"error": "function service unreachable"})
		return
	}

	defer resp.Body.Close()

	GinContext.Status(resp.StatusCode)
	for k, vv := range resp.Header {
		GinContext.Writer.Header()[k] = vv
	}

	io.Copy(GinContext.Writer, resp.Body)
}
