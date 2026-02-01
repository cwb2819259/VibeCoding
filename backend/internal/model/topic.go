package model

import "time"

type Topic struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:64;not null;default:'';uniqueIndex" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Topic) TableName() string { return "topics" }

type PostTopic struct {
	PostID    uint64    `gorm:"primaryKey" json:"post_id"`
	TopicID   uint64    `gorm:"primaryKey;index" json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PostTopic) TableName() string { return "post_topics" }
