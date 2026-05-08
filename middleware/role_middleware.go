package middleware

import (
	"go-auth-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil role user yang sudah disimpan oleh AuthMiddleware sebelumnya
		userRole, exists := c.Get("currentUser_Role")
		if !exists {
			utils.JSONResponse(c, http.StatusForbidden, "error", "Role tidak ditemukan", nil)
			c.Abort()
			return
		}

		// Cek apakah role user ada dalam daftar role yang diizinkan
		isAllowed := false
		for _, role := range roles {
			if role == userRole.(string) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			utils.JSONResponse(c, http.StatusForbidden, "error", "Anda tidak memiliki akses ke halaman ini", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
