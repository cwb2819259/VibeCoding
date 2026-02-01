package model

import "time"

type Comment struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint64    `gorm:"not null;default:0;index" json:"post_id"`
	UserID    uint64    `gorm:"not null;default:0;index" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User *User `gorm:"-" json:"user,omitempty"`
}

func (Comment) TableName() string { return "comments" }
