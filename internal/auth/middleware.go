package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(Context *gin.Context) {
		authHeader := Context.GetHeader("Authorization")

		if !hasAuthorizationHeader(authHeader, Context) {
			return
		}

		tokenString, ok := extractBearerToken(authHeader, Context)
		if !ok {
			return
		}

		claims := &Claims{}
		if !isTokenValid(tokenString, claims, Context) {
			return
		}

		Context.Set("username", claims.Username)
		Context.Set("role", claims.Role)

		Context.Next()
	}
}

func hasAuthorizationHeader(header string, context *gin.Context) bool {
	if header == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header eksik"})
		return false
	}
	return true
}

func extractBearerToken(header string, context *gin.Context) (string, bool) {
	fields := strings.Fields(header)
	if len(fields) != 2 || fields[0] != "Bearer" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header formatı yanlış"})
		return "", false
	}
	return fields[1], true
}

func isTokenValid(tokenString string, claims *Claims, context *gin.Context) bool {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz veya süresi dolmuş token"})
		return false
	}
	return true
}
