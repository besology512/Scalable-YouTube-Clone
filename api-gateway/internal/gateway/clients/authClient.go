package clients

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthClient struct {
	baseURL string
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{baseURL: baseURL}
}

func (a *AuthClient) GetBaseURL() string {
	return a.baseURL
}
func (request *AuthClient) ForwardAuth(GinContext *gin.Context) {
	target := request.baseURL + GinContext.Request.URL.Path

	fmt.Println(target)

	req, _ := http.NewRequestWithContext(GinContext.Request.Context(), GinContext.Request.Method, target, GinContext.Request.Body)

	for k, v := range GinContext.Request.Header {
		req.Header[k] = v
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		GinContext.JSON(http.StatusServiceUnavailable, gin.H{"error": "Auth service unreachable"})
		return
	}
	defer resp.Body.Close()

	GinContext.Status(resp.StatusCode)
	for k, vv := range resp.Header {
		GinContext.Writer.Header()[k] = vv
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		GinContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}
	GinContext.Writer.Write(body)
}

func (request *AuthClient) LogOut(GinContext *gin.Context) {
	request.ForwardAuth(GinContext)
}
func (request *AuthClient) Login(GinContext *gin.Context) {
	request.ForwardAuth(GinContext)
}
func (request *AuthClient) Refresh(GinContext *gin.Context) {
	request.ForwardAuth(GinContext)
}
func (request *AuthClient) Health(GinContext *gin.Context) {
	request.ForwardAuth(GinContext)
}
