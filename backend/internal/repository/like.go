package repository

import (
	"github.com/vibecoding/community/internal/model"
	"gorm.io/gorm"
)

type LikeRepo struct {
	db *gorm.DB
}

func NewLikeRepo(db *gorm.DB) *LikeRepo { return &LikeRepo{db: db} }

func (r *LikeRepo) Create(l *model.Like) error {
	return r.db.Create(l).Error
}

func (r *LikeRepo) Delete(userID, targetID uint64, targetType model.TargetType) error {
	return r.db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Delete(&model.Like{}).Error
}

func (r *LikeRepo) Exists(userID, targetID uint64, targetType model.TargetType) (bool, error) {
	var c int64
	err := r.db.Model(&model.Like{}).Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Count(&c).Error
	return c > 0, err
}

// CountByTarget 按目标隔离统计：只统计 target_type + target_id 匹配的点赞数（帖子维度）
func (r *LikeRepo) CountByTarget(targetID uint64, targetType model.TargetType) (int64, error) {
	var c int64
	err := r.db.Model(&model.Like{}).Where("target_type = ? AND target_id = ?", targetType, targetID).Count(&c).Error
	return c, err
}
