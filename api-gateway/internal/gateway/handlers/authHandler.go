package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/besology512/api-gateway/internal/gateway/clients"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	client *clients.AuthClient
}

func SetAuthClient(client *clients.AuthClient) *AuthHandler {
	return &AuthHandler{client: client}
}

func (request *AuthHandler) Proxy(GinContext *gin.Context) {

	target := request.client.GetBaseURL() + GinContext.Request.URL.Path
	req, err := http.NewRequestWithContext(GinContext.Request.Context(), GinContext.Request.Method, target, GinContext.Request.Body)

	fmt.Println(target)

	if err != nil {
		GinContext.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

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

	io.Copy(GinContext.Writer, resp.Body)
}

func (request *AuthHandler) LogOut(GinContext *gin.Context) {

	//refresh token from cookie and put in body
	// access token in header (use middleware)
	//delete cookie from user (httpOnly)
	request.Proxy(GinContext)
}
func (request *AuthHandler) Login(GinContext *gin.Context) {
	//does nothing but redirect to callback
	request.Proxy(GinContext)
}
func (request *AuthHandler) Refresh(GinContext *gin.Context) {
	//extract refresh token from cookie
	//send in body
	request.Proxy(GinContext)
}
func (request *AuthHandler) Health(GinContext *gin.Context) {
	request.Proxy(GinContext)
}
func (request *AuthHandler) Callback(GinContext *gin.Context) {
	//redirect refresh and access token
	//access in json body
	//refresh in cookie

	request.Proxy(GinContext)
}
