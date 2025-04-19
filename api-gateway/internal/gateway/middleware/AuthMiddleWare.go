package middleware

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(publicKey *rsa.PublicKey) gin.HandlerFunc {

	return func(GinContext *gin.Context) {
		header := GinContext.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			GinContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", jwtToken.Header["alg"])
			}
			return publicKey, nil
		})

		if err != nil || !token.Valid {
			GinContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			GinContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
		GinContext.Set("user", claims)
		GinContext.Next()
	}

}
