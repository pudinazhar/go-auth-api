package repository

import (
	"go-auth-api/model"

	"gorm.io/gorm"
)

type TokenRepository interface {
	SaveToken(token model.RefreshToken) error
	DeleteToken(token string) error
	FindByToken(token string) (model.RefreshToken, error)
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *tokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) SaveToken(token model.RefreshToken) error {
	return r.db.Create(&token).Error
}

func (r *tokenRepository) DeleteToken(token string) error {
	// Proses revoke token saat logout
	return r.db.Where("token = ?", token).Delete(&model.RefreshToken{}).Error
}

func (r *tokenRepository) FindByToken(token string) (model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.Where("token = ?", token).First(&rt).Error
	return rt, err
}
