package repository

import (
	"github.com/vibecoding/community/internal/model"
	"gorm.io/gorm"
)

type PostRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) *PostRepo { return &PostRepo{db: db} }

func (r *PostRepo) Create(p *model.Post) error {
	return r.db.Create(p).Error
}

func (r *PostRepo) GetByID(id uint64) (*model.Post, error) {
	var p model.Post
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// List 列表。status 为空或 "all" 时不按状态过滤（管理端「全部」）；为空时 C 端默认只查 normal。
func (r *PostRepo) List(offset, limit int, status model.PostStatus, keyword string, topicID uint64) ([]model.Post, int64, error) {
	q := r.db.Model(&model.Post{})
	s := string(status)
	// 仅在有明确状态时加筛选：normal / hidden / flagged
	if s == string(model.PostStatusNormal) || s == string(model.PostStatusHidden) || s == string(model.PostStatusFlagged) {
		q = q.Where("status = ?", status)
	}
	// status 为空或 "all" 时不加 status 条件，返回所有状态
	if keyword != "" {
		q = q.Where("content LIKE ?", "%"+keyword+"%")
	}
	if topicID > 0 {
		q = q.Joins("JOIN post_topics ON post_topics.post_id = posts.id AND post_topics.topic_id = ?", topicID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Post
	err := q.Order("posts.created_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *PostRepo) ListByUserID(userID uint64, offset, limit int) ([]model.Post, int64, error) {
	q := r.db.Model(&model.Post{}).Where("user_id = ? AND status = ?", userID, model.PostStatusNormal)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Post
	err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *PostRepo) Update(p *model.Post) error {
	return r.db.Save(p).Error
}

func (r *PostRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Post{}, id).Error
}

func (r *PostRepo) GetMediaByPostID(postID uint64) ([]model.PostMedia, error) {
	var list []model.PostMedia
	err := r.db.Where("post_id = ?", postID).Order("sort_order ASC").Find(&list).Error
	return list, err
}

func (r *PostRepo) CreateMedia(m *model.PostMedia) error {
	return r.db.Create(m).Error
}

func (r *PostRepo) DeleteMediaByPostID(postID uint64) error {
	return r.db.Where("post_id = ?", postID).Delete(&model.PostMedia{}).Error
}
