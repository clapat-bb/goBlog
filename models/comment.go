package models

import "time"

type Conment struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Content string `gorm:"type:text;not null" json:"content"`
	UserID  uint   `json:"user_id"`
	User    User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`

	PostID uint `json:"post_id"`
	Post   Post `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	ParentID *uint     `json:"parent_id"`
	Replies  []Conment `gorm:"foreignKey:ParentID" json:"replies"`

	CreatedAT time.Time `json:"created_at"`
	UpdatedAt time.Time `josn:"updated_at"`
	DeletedAt time.Time `gorm:"index" json:"-"`
}
