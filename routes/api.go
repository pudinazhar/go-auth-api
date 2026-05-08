package routes

import (
	"go-auth-api/controller"
	_ "go-auth-api/docs"
	"go-auth-api/middleware"

	"github.com/gin-gonic/gin"
	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter mengatur semua rute API
func SetupRouter(r *gin.Engine, authController *controller.AuthController) {

	// Endpoint Dokumentasi Swagger
	// Gunakan swagFiles.Handler, bukan swaggerFiles.Handler
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swagFiles.Handler))

	api := r.Group("/api/v1")
	{
		// Public Routes
		api.POST("/register", authController.Register)
		api.POST("/login", authController.Login)
		api.POST("/refresh", authController.RefreshToken)
		api.POST("/forgot-password", authController.ForgotPassword)
		api.POST("/reset-password", authController.ResetPassword)
		// Protected Routes (Harus Login)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", authController.GetProfile)
			protected.POST("/logout", authController.Logout)

			// Khusus Admin
			adminOnly := protected.Group("/admin")
			adminOnly.Use(middleware.RoleMiddleware("Admin"))
			{
				adminOnly.GET("/dashboard", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Welcome Admin!"})
				})
			}
		}
	}
}
