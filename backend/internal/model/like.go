package model

import "time"

type TargetType string

const (
	TargetTypePost    TargetType = "post"
	TargetTypeComment TargetType = "comment"
)

type Like struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint64     `gorm:"not null;default:0;uniqueIndex:uk_user_target" json:"user_id"`
	TargetType TargetType `gorm:"type:enum('post','comment');not null;default:'post'" json:"target_type"`
	TargetID   uint64     `gorm:"not null;default:0;uniqueIndex:uk_user_target" json:"target_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (Like) TableName() string { return "likes" }
