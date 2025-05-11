package handler

import (
	"auth-service/internal/auth"
	"auth-service/internal/service"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: s}
}

// GoogleLogin godoc
// @Summary Google OAuth login
// @Description Redirect user to Google's OAuth login page
// @Tags Auth
// @Produce json
// @Success 302 {string} string "Redirect to Google"
// @Router /auth/google/login [get]

// GoogleLogin redirects user to Google OAuth page
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	c.Request = c.Request.WithContext(context.WithValue(
		c.Request.Context(),
		"provider",
		"google",
	))
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// GoogleCallback godoc
// @Summary Google OAuth callback
// @Description Handles Google callback, registers user, and returns tokens
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/google/callback [get]

// GoogleCallback handles the callback from Google and returns token pair
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Gothic OAuth2 handling
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed: " + err.Error()})
		return
	}

	// Business logic (Persisting user & Generating tokens)
	tokens, err := h.Service.HandleGoogleCallback(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout godoc
// @Summary Logout user
// @Description Revoke refresh token by its ID (JTI)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body struct{RefreshToken string `json:"refresh_token"`} true "Refresh Token"
// @Success 200 {string} string "OK"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]

// Logout revokes a refresh token using the token's JTI
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	claims, err := auth.VerifyToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	if err := h.Service.Logout(claims.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke token"})
		return
	}

	c.Status(http.StatusOK)
}

// Refresh godoc
// @Summary Refresh tokens
// @Description Generate new access and refresh tokens from a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body struct{RefreshToken string `json:"refresh_token"`} true "Refresh Token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]

// Refresh generates a new access/refresh token pair from a valid refresh token
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	tokens, err := h.Service.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
