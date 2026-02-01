package model

import "time"

type PostType string
type PostStatus string

const (
	PostTypeText  PostType = "text"
	PostTypeImage PostType = "image"
	PostTypeVideo PostType = "video"
)

const (
	PostStatusNormal  PostStatus = "normal"
	PostStatusHidden  PostStatus = "hidden"
	PostStatusFlagged PostStatus = "flagged"
)

type Post struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64     `gorm:"not null;default:0;index" json:"user_id"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	Type      PostType   `gorm:"type:enum('text','image','video');not null;default:'text'" json:"type"`
	Status    PostStatus `gorm:"type:enum('normal','hidden','flagged');not null;default:'normal'" json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	User   *User       `gorm:"-" json:"user,omitempty"`
	Media  []PostMedia `gorm:"-" json:"media,omitempty"`
}

func (Post) TableName() string { return "posts" }

type PostMedia struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint64    `gorm:"not null;default:0;index" json:"post_id"`
	Type      string    `gorm:"type:enum('image','video');not null;default:'image'" json:"type"`
	URL       string    `gorm:"size:512;not null;default:''" json:"url"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PostMedia) TableName() string { return "post_media" }
