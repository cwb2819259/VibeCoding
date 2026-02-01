package service

import (
	"time"

	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/internal/repository"
)

// unreadSentinel 未读标记时间，与表默认一致，避免 GORM 零值 0001-01-01 超出 MySQL DATETIME 范围导致插入失败
var unreadSentinel = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

type NotificationService struct {
	repo *repository.NotificationRepo
}

func NewNotificationService(repo *repository.NotificationRepo) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Create(userID uint64, typ string, relatedID uint64, payload model.JSONMap) error {
	n := &model.Notification{
		UserID:    userID,
		Type:      typ,
		RelatedID: relatedID,
		Payload:   payload,
		ReadAt:    unreadSentinel, // 显式未读，避免 GORM 零值 0001-01-01 导致 MySQL 插入失败
	}
	return s.repo.Create(n)
}

func (s *NotificationService) ListByUserID(userID uint64, offset, limit int) ([]model.Notification, int64, error) {
	return s.repo.ListByUserID(userID, offset, limit)
}

func (s *NotificationService) MarkRead(id, userID uint64) error {
	return s.repo.MarkRead(id, userID)
}

func (s *NotificationService) MarkAllRead(userID uint64) error {
	return s.repo.MarkAllRead(userID)
}

func (s *NotificationService) CountUnread(userID uint64) (int64, error) {
	return s.repo.CountUnreadByUserID(userID)
}
