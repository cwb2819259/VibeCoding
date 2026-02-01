package model

import "time"

type Admin struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string    `gorm:"size:255;not null;default:''" json:"-"`
	Nickname     string    `gorm:"size:64;not null;default:''" json:"nickname"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Admin) TableName() string { return "admins" }
