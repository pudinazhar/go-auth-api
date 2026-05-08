package service

import (
	"errors"
	"go-auth-api/model"
	"go-auth-api/repository"
	"go-auth-api/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Register(input model.RegisterInput) (model.User, error)
	Login(input model.LoginInput) (string, string, model.User, error)
	GetByID(id uint) (model.User, error)
	Logout(token string) error
	RefreshToken(tokenStr string) (string, string, error)
	ForgotPassword(email string) error
	ResetPassword(token string, newPassword string) error
}

type authService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthService(ur repository.UserRepository, tr repository.TokenRepository) *authService {
	return &authService{ur, tr}
}

func (s *authService) Register(input model.RegisterInput) (model.User, error) {
	// 1. Cek apakah email sudah terdaftar
	_, err := s.userRepo.FindByEmail(input.Email)
	if err == nil {
		return model.User{}, errors.New("email sudah digunakan")
	}

	// 2. Hash Password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return model.User{}, err
	}

	// 3. Simpan ke database
	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: hashedPassword,
		Role:     input.Role,
	}

	return s.userRepo.Save(user)
}

func (s *authService) Login(input model.LoginInput) (string, string, model.User, error) {
	// 1. Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return "", "", model.User{}, errors.New("email atau password salah")
	}

	// 2. Verifikasi Password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return "", "", model.User{}, errors.New("email atau password salah")
	}

	// 3. Generate Access Token (15 Menit)
	accessToken, _ := utils.GenerateToken(user.ID, user.Role, 15*time.Minute)

	// 4. Generate Refresh Token (7 Hari)
	refreshTokenStr, _ := utils.GenerateToken(user.ID, user.Role, 7*24*time.Hour)

	// 5. Simpan Refresh Token ke Database (untuk fitur logout/revoke nanti)
	refreshTokenRecord := model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour),
	}
	s.tokenRepo.SaveToken(refreshTokenRecord)

	return accessToken, refreshTokenStr, user, nil
}

// 2. Tambahkan Implementasi fungsinya di bawah
func (s *authService) GetByID(id uint) (model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// LogOut: Hapus Refresh Token dari Database
func (s *authService) Logout(token string) error {
	return s.tokenRepo.DeleteToken(token)
}

// RefreshToken: Validasi refresh token, buat access token baru, dan (opsional) buat refresh token baru
func (s *authService) RefreshToken(tokenStr string) (string, string, error) {
	// A. Cari token di database (pastikan token belum di-revoke/logout)
	storedToken, err := s.tokenRepo.FindByToken(tokenStr)
	if err != nil {
		return "", "", errors.New("refresh token tidak valid atau sudah logout")
	}

	// B. Cek apakah token di database sudah expired secara waktu
	if time.Now().After(storedToken.ExpiredAt) {
		s.tokenRepo.DeleteToken(tokenStr) // Hapus yang sudah basi
		return "", "", errors.New("refresh token sudah expired, silahkan login ulang")
	}

	// C. Validasi JWT nya (untuk mengambil UserID & Role)
	token, err := utils.ValidateToken(tokenStr)
	if err != nil || !token.Valid {
		return "", "", errors.New("token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("gagal memproses data token")
	}

	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	// D. Buat Access Token Baru (15 Menit)
	newAccessToken, _ := utils.GenerateToken(userID, role, 15*time.Minute)

	// E. (Opsional) Buat Refresh Token Baru agar masa berlakunya bertambah (Rotate)
	newRefreshToken, _ := utils.GenerateToken(userID, role, 7*24*time.Hour)

	// Update token di database (Hapus yang lama, simpan yang baru)
	s.tokenRepo.DeleteToken(tokenStr)
	s.tokenRepo.SaveToken(model.RefreshToken{
		UserID:    userID,
		Token:     newRefreshToken,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour),
	})

	return newAccessToken, newRefreshToken, nil
}

// Tambahkan ke Interface AuthServic
func (s *authService) ForgotPassword(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("email tidak ditemukan")
	}

	// Generate Token Unik (Bisa pakai UUID atau string acak)
	resetToken := utils.GenerateRandomString(32)
	user.ResetPasswordToken = resetToken

	_, err = s.userRepo.Update(user)
	return err
}

func (s *authService) ResetPassword(token string, newPassword string) error {
	user, err := s.userRepo.FindByResetToken(token)
	if err != nil {
		return errors.New("token reset tidak valid")
	}

	// Hash password baru
	hashedPassword, _ := utils.HashPassword(newPassword)
	user.Password = hashedPassword
	user.ResetPasswordToken = "" // Hapus token setelah dipakai

	_, err = s.userRepo.Update(user)
	return err
}
