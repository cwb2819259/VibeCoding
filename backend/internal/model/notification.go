package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		*j = nil
		return nil
	}
	return json.Unmarshal(b, j)
}

type Notification struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null;default:0;index" json:"user_id"`
	Type      string    `gorm:"size:32;not null;default:''" json:"type"`
	RelatedID uint64    `gorm:"not null;default:0" json:"related_id"`
	Payload   JSONMap   `gorm:"type:json" json:"payload,omitempty"`
	ReadAt    time.Time `gorm:"not null" json:"read_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Notification) TableName() string { return "notifications" }
