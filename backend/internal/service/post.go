package service

import (
	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/internal/repository"
)

type PostService struct {
	postRepo   *repository.PostRepo
	userRepo   *repository.UserRepo
	topicRepo  *repository.TopicRepo
	likeRepo   *repository.LikeRepo
	commentRepo *repository.CommentRepo
}

func NewPostService(
	postRepo *repository.PostRepo,
	userRepo *repository.UserRepo,
	topicRepo *repository.TopicRepo,
	likeRepo *repository.LikeRepo,
	commentRepo *repository.CommentRepo,
) *PostService {
	return &PostService{
		postRepo:   postRepo,
		userRepo:   userRepo,
		topicRepo:  topicRepo,
		likeRepo:   likeRepo,
		commentRepo: commentRepo,
	}
}

func (s *PostService) Create(userID uint64, content string, postType model.PostType, mediaURLs []string, topicNames []string) (*model.Post, error) {
	p := &model.Post{
		UserID:  userID,
		Content: content,
		Type:    postType,
		Status:  model.PostStatusNormal,
	}
	if err := s.postRepo.Create(p); err != nil {
		return nil, err
	}
	for i, url := range mediaURLs {
		t := "image"
		if postType == model.PostTypeVideo {
			t = "video"
		}
		m := &model.PostMedia{PostID: p.ID, Type: t, URL: url, SortOrder: i}
		_ = s.postRepo.CreateMedia(m)
	}
	for _, name := range topicNames {
		if name == "" {
			continue
		}
		topic, err := s.topicRepo.GetOrCreateByName(name)
		if err != nil {
			continue
		}
		_ = s.topicRepo.AddPostTopic(p.ID, topic.ID)
	}
	return p, nil
}

func (s *PostService) GetByID(id uint64, loadUser, loadMedia bool) (*model.Post, error) {
	p, err := s.postRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if loadUser && p.UserID > 0 {
		u, _ := s.userRepo.GetByID(p.UserID)
		p.User = u
	}
	if loadMedia {
		media, _ := s.postRepo.GetMediaByPostID(p.ID)
		p.Media = media
	}
	return p, nil
}

func (s *PostService) List(offset, limit int, status model.PostStatus, keyword string, topicID uint64) ([]model.Post, int64, error) {
	list, total, err := s.postRepo.List(offset, limit, status, keyword, topicID)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		u, _ := s.userRepo.GetByID(list[i].UserID)
		list[i].User = u
		media, _ := s.postRepo.GetMediaByPostID(list[i].ID)
		list[i].Media = media
	}
	return list, total, err
}

func (s *PostService) ListByUserID(userID uint64, offset, limit int) ([]model.Post, int64, error) {
	return s.postRepo.ListByUserID(userID, offset, limit)
}

// ListByUserIDWithDetails 我的帖子列表，附带 user、media（与帖子列表展示一致）
func (s *PostService) ListByUserIDWithDetails(userID uint64, offset, limit int) ([]model.Post, int64, error) {
	list, total, err := s.postRepo.ListByUserID(userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		u, _ := s.userRepo.GetByID(list[i].UserID)
		list[i].User = u
		media, _ := s.postRepo.GetMediaByPostID(list[i].ID)
		list[i].Media = media
	}
	return list, total, nil
}

func (s *PostService) UpdateStatus(id uint64, status model.PostStatus) error {
	p, err := s.postRepo.GetByID(id)
	if err != nil {
		return err
	}
	p.Status = status
	return s.postRepo.Update(p)
}

func (s *PostService) Delete(id uint64) error {
	_ = s.postRepo.DeleteMediaByPostID(id)
	return s.postRepo.Delete(id)
}

func (s *PostService) LikeCount(postID uint64) (int64, error) {
	return s.likeRepo.CountByTarget(postID, model.TargetTypePost)
}

func (s *PostService) HasLiked(userID, postID uint64) (bool, error) {
	return s.likeRepo.Exists(userID, postID, model.TargetTypePost)
}

func (s *PostService) CommentCount(postID uint64) (int64, error) {
	return s.commentRepo.CountByPostID(postID)
}
