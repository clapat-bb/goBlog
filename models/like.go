package models

import (
	"time"

	"gorm.io/gorm"
)

type Like struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `json:"user_id"`
	TargetID   uint           `json:"target_id"`
	TargetType string         `json:"target_type"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
