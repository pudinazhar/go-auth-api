package main

import (
	"go-auth-api/config"
	"go-auth-api/controller"
	"go-auth-api/model"
	"go-auth-api/repository"
	"go-auth-api/routes"
	"go-auth-api/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// @title           Go Auth API
// @version         1.0
// @description     API Documentation for Golang Authentication System.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Inisialisasi Database
	config.ConnectDatabase()
	err := config.DB.AutoMigrate(&model.User{}, &model.RefreshToken{}, &model.LoginLog{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	// 2. Inisialisasi Dependency (Repository -> Service -> Controller)
	userRepo := repository.NewUserRepository(config.DB)
	tokenRepo := repository.NewTokenRepository(config.DB)
	logRepo := repository.NewLogRepository(config.DB)

	authService := service.NewAuthService(userRepo, tokenRepo)
	authController := controller.NewAuthController(authService, logRepo)

	// 3. Setup Framework & Routes
	r := gin.Default()

	// Panggil fungsi route yang kita buat di folder routes
	routes.SetupRouter(r, authController)

	// 4. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
