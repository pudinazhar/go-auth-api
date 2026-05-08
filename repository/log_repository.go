package repository

import (
	"go-auth-api/model"

	"gorm.io/gorm"
)

type LogRepository interface {
	SaveLog(log model.LoginLog) error
}

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) *logRepository {
	return &logRepository{db}
}

func (r *logRepository) SaveLog(log model.LoginLog) error {
	return r.db.Create(&log).Error
}
