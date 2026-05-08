package middleware

import (
	"go-auth-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			utils.JSONResponse(c, http.StatusUnauthorized, "error", "Silahkan login terlebih dahulu", nil)
			c.Abort()
			return
		}

		// 2. Ekstrak token-nya saja
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 3. Validasi Token
		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			utils.JSONResponse(c, http.StatusUnauthorized, "error", "Token tidak valid atau expired", nil)
			c.Abort()
			return
		}

		// 4. Ambil data dari claims dan simpan di context Gin
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			c.Set("currentUser_ID", uint(claims["user_id"].(float64)))
			c.Set("currentUser_Role", claims["role"].(string))
		}

		c.Next()
	}
}
