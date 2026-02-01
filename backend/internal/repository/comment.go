package repository

import (
	"github.com/vibecoding/community/internal/model"
	"gorm.io/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo { return &CommentRepo{db: db} }

func (r *CommentRepo) Create(c *model.Comment) error {
	return r.db.Create(c).Error
}

func (r *CommentRepo) ListByPostID(postID uint64, offset, limit int) ([]model.Comment, int64, error) {
	q := r.db.Model(&model.Comment{}).Where("post_id = ?", postID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Comment
	err := q.Order("created_at ASC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

// CountByPostID 按帖子隔离统计：只统计该 post_id 下的评论数（帖子维度）
func (r *CommentRepo) CountByPostID(postID uint64) (int64, error) {
	var c int64
	err := r.db.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&c).Error
	return c, err
}
