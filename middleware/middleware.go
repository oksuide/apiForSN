package middleware

import (
	"apiForSN/auth"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Middleware для проверки JWT токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Токен должен начинаться с "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is missing"})
			c.Abort()
			return
		}

		// Парсим и проверяем токен
		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Сохраняем userID в контексте для дальнейшего использования
		c.Set("userID", claims.UserID)

		// Переходим к следующему обработчику
		c.Next()
	}
}

func PostIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("postID")
		if postID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "postID is required"})
			c.Abort()
			return
		}

		// Конвертируем postID в int
		id, err := strconv.Atoi(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid postID"})
			c.Abort()
			return
		}

		// Сохраняем postID в контексте
		c.Set("postID", id)
		c.Next()
	}
}

func CommentIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		commentID := c.Param("commentID")
		if commentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "commentID is required"})
			c.Abort()
			return
		}

		// Конвертируем commentID в int
		id, err := strconv.Atoi(commentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid commentID"})
			c.Abort()
			return
		}

		// Сохраняем commentID в контексте
		c.Set("commentID", id)
		c.Next()
	}
}
