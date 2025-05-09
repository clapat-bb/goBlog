package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `grom:"unique;not null" json:"name"`
	Posts     []*Post        `gorm:"many2many:post_tags;" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeleteAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
