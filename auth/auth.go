package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Секретный ключ для подписи токенов. Его нужно хранить в защищенном месте (например, в переменной окружения).
var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Структура данных для токена
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Функция генерации токена
func GenerateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Устанавливаем срок действия токена на 24 часа
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
