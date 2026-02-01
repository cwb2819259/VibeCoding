package repository

import (
	"time"

	"github.com/vibecoding/community/internal/model"
	"gorm.io/gorm"
)

type NotificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) *NotificationRepo { return &NotificationRepo{db: db} }

func (r *NotificationRepo) Create(n *model.Notification) error {
	return r.db.Create(n).Error
}

func (r *NotificationRepo) ListByUserID(userID uint64, offset, limit int) ([]model.Notification, int64, error) {
	q := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Notification
	err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *NotificationRepo) MarkRead(id, userID uint64) error {
	return r.db.Model(&model.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("read_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

func (r *NotificationRepo) MarkAllRead(userID uint64) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ?", userID).Update("read_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}

// CountUnreadByUserID 未读数量（read_at 早于 1970-01-02 视为未读，兼容零值与表默认）
func (r *NotificationRepo) CountUnreadByUserID(userID uint64) (int64, error) {
	var c int64
	threshold := time.Date(1970, 1, 2, 0, 0, 0, 0, time.UTC)
	err := r.db.Model(&model.Notification{}).Where("user_id = ? AND read_at < ?", userID, threshold).Count(&c).Error
	return c, err
}
