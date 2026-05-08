package controller

import (
	"go-auth-api/model"
	"go-auth-api/repository"
	"go-auth-api/service"
	"go-auth-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
	logRepo     repository.LogRepository // Kita inject repository log di sini
}

func NewAuthController(as service.AuthService, lr repository.LogRepository) *AuthController {
	return &AuthController{as, lr}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with name, email, phone, password, and role
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      model.RegisterInput  true  "Register User"
// @Success      201      {object}  utils.Response
// @Failure      400      {object}  utils.Response
// @Router       /register [post]
func (h *AuthController) Register(c *gin.Context) {
	var input model.RegisterInput

	// 1. Validasi Input Body
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	// 2. Panggil Service Register
	user, err := h.authService.Register(input)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "success", "Registrasi berhasil", user)
}

// @Summary      Login user
// @Description  Authenticate user and return tokens
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      model.LoginInput  true  "Login User"
// @Success      200      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Router       /login [post]
func (h *AuthController) Login(c *gin.Context) {
	var input model.LoginInput

	// 1. Validasi Input
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	// 2. Panggil Service Login
	accessToken, refreshToken, user, err := h.authService.Login(input)

	// Siapkan data untuk Log
	logEntry := model.LoginLog{
		UserID:    user.ID,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	if err != nil {
		logEntry.Status = "FAILED"
		h.logRepo.SaveLog(logEntry) // Simpan log gagal
		utils.JSONResponse(c, http.StatusUnauthorized, "error", err.Error(), nil)
		return
	}

	// 3. Simpan Log Sukses
	logEntry.Status = "SUCCESS"
	h.logRepo.SaveLog(logEntry)

	// 4. Kirim Response
	utils.JSONResponse(c, http.StatusOK, "success", "Login berhasil", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// @Summary      Get user profile
// @Description  Get current logged in user profile information
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  utils.Response
// @Failure      401      {object}  utils.Response
// @Router       /profile [get]
func (h *AuthController) GetProfile(c *gin.Context) {
	// Ambil ID dari context (hasil set di middleware)
	userID, _ := c.Get("currentUser_ID")

	// Panggil service/repo untuk cari user berdasarkan ID
	// Kita gunakan userRepo langsung atau lewat service
	user, err := h.authService.GetByID(userID.(uint))
	if err != nil {
		utils.JSONResponse(c, http.StatusNotFound, "error", "User tidak ditemukan", nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "success", "Data profil berhasil diambil", user)
}

// @Summary      Refresh Access Token
// @Description  Get a new access token using a valid refresh token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Refresh Token"
// @Success      200      {object}  utils.Response
// @Router       /refresh [post]
func (h *AuthController) RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// 1. Validasi Input
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "error", "Refresh token wajib diisi", nil)
		return
	}

	// 2. Panggil Service
	newAccess, newRefresh, err := h.authService.RefreshToken(input.RefreshToken)
	if err != nil {
		utils.JSONResponse(c, http.StatusUnauthorized, "error", err.Error(), nil)
		return
	}

	// 3. Berikan response token baru
	utils.JSONResponse(c, http.StatusOK, "success", "Token berhasil diperbarui", gin.H{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
	})
}

// @Summary      Logout user
// @Description  Revoke refresh token and logout
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      object  true  "Refresh Token"
// @Success      200      {object}  utils.Response
// @Router       /logout [post]
func (h *AuthController) Logout(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, "error", "Refresh token diperlukan", nil)
		return
	}

	err := h.authService.Logout(input.RefreshToken)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, "error", "Gagal logout", nil)
		return
	}

	utils.JSONResponse(c, http.StatusOK, "success", "Logout berhasil", nil)
}

// @Summary      Forgot Password
// @Description  Request a password reset token via email
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "User Email"
// @Success      200      {object}  utils.Response
// @Router       /forgot-password [post]
func (h *AuthController) ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.JSONResponse(c, 400, "error", "Email tidak valid", nil)
		return
	}

	err := h.authService.ForgotPassword(input.Email)
	if err != nil {
		// Demi keamanan, biasanya kita tetap kasih response success
		// agar penyerang tidak tahu email mana yang terdaftar.
		// Tapi untuk belajar, kita tampilkan errornya.
		utils.JSONResponse(c, 404, "error", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, 200, "success", "Token reset telah dibuat (cek database kolom reset_password_token)", nil)
}

// @Summary      Reset Password
// @Description  Reset password using the reset token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      object  true  "Token and New Password"
// @Success      200      {object}  utils.Response
// @Router       /reset-password [post]
func (h *AuthController) ResetPassword(c *gin.Context) {
	var input struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.JSONResponse(c, 400, "error", err.Error(), nil)
		return
	}

	err := h.authService.ResetPassword(input.Token, input.Password)
	if err != nil {
		utils.JSONResponse(c, 400, "error", err.Error(), nil)
		return
	}

	utils.JSONResponse(c, 200, "success", "Password berhasil diubah, silakan login kembali", nil)
}
