package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims adalah struktur data yang kita titipkan di dalam token
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.MapClaims
}

// GenerateToken membuat JWT token (bisa digunakan untuk Access atau Refresh Token)
func GenerateToken(userID uint, role string, duration time.Duration) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ValidateToken memeriksa apakah token asli dan belum expired
func ValidateToken(tokenString string) (*jwt.Token, error) {
	secretKey := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode signing-nya HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
}
