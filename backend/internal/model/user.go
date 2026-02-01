package model

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone     string    `gorm:"size:20;not null;uniqueIndex" json:"phone"`
	Nickname  string    `gorm:"size:64;not null;default:''" json:"nickname"`
	AvatarURL string    `gorm:"size:512;not null;default:''" json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }
