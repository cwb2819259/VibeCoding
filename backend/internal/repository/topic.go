package repository

import (
	"github.com/vibecoding/community/internal/model"
	"gorm.io/gorm"
)

type TopicRepo struct {
	db *gorm.DB
}

func NewTopicRepo(db *gorm.DB) *TopicRepo { return &TopicRepo{db: db} }

func (r *TopicRepo) GetOrCreateByName(name string) (*model.Topic, error) {
	var t model.Topic
	err := r.db.Where("name = ?", name).First(&t).Error
	if err == nil {
		return &t, nil
	}
	t.Name = name
	if err := r.db.Create(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TopicRepo) ListAll() ([]model.Topic, error) {
	var list []model.Topic
	err := r.db.Order("name ASC").Find(&list).Error
	return list, err
}

func (r *TopicRepo) AddPostTopic(postID, topicID uint64) error {
	pt := model.PostTopic{PostID: postID, TopicID: topicID}
	return r.db.Create(&pt).Error
}

func (r *TopicRepo) GetTopicIDsByPostID(postID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.db.Model(&model.PostTopic{}).Where("post_id = ?", postID).Pluck("topic_id", &ids).Error
	return ids, err
}
