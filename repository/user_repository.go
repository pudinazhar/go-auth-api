package repository

import (
	"go-auth-api/model"

	"gorm.io/gorm"
)

// UserRepository adalah interface (kontrak) yang mendefinisikan fungsi apa saja yang tersedia
type UserRepository interface {
	Save(user model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindByID(id uint) (model.User, error)
	Update(user model.User) (model.User, error)
	FindByResetToken(token string) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository adalah constructor untuk inisialisasi repository
func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Save(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *userRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) FindByID(id uint) (model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *userRepository) Update(user model.User) (model.User, error) {
	err := r.db.Save(&user).Error
	return user, err
}

func (r *userRepository) FindByResetToken(token string) (model.User, error) {
	var user model.User
	err := r.db.Where("reset_password_token = ?", token).First(&user).Error
	return user, err
}
