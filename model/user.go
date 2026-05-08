package model

import (
	"time"
)

type User struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	Name               string    `gorm:"type:varchar(100);not null" json:"name"`
	Email              string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Phone              string    `gorm:"type:varchar(20);not null" json:"phone"`
	Password           string    `gorm:"type:varchar(255);not null" json:"-"` // "-" artinya tidak muncul di JSON
	Role               string    `gorm:"type:enum('Admin', 'Guest');default:'Guest'" json:"role"`
	ResetPasswordToken string    `gorm:"type:varchar(255)" json:"-"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Token     string    `gorm:"type:text;not null"`
	ExpiredAt time.Time `json:"expired_at"`
}

type LoginLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=Admin Guest"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
