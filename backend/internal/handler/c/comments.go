package c

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vibecoding/community/internal/middleware"
	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/internal/notification"
	"github.com/vibecoding/community/internal/repository"
)

type CommentsHandler struct {
	commentRepo *repository.CommentRepo
	postRepo    *repository.PostRepo
	userRepo    *repository.UserRepo
	notifProd   *notification.Producer
}

func NewCommentsHandler(commentRepo *repository.CommentRepo, postRepo *repository.PostRepo, userRepo *repository.UserRepo, notifProd *notification.Producer) *CommentsHandler {
	return &CommentsHandler{commentRepo: commentRepo, postRepo: postRepo, userRepo: userRepo, notifProd: notifProd}
}

// List 评论列表
// @Summary 评论列表
// @Tags C端-互动
// @Produce json
// @Param id path int true "帖子ID"
// @Param offset query int false "0"
// @Param limit query int false "20"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/posts/{id}/comments [get]
func (h *CommentsHandler) List(c *gin.Context) {
	postID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if postID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	list, total, err := h.commentRepo.ListByPostID(postID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range list {
		u, _ := h.userRepo.GetByID(list[i].UserID)
		list[i].User = u
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// Create 发表评论
// @Summary 发表评论
// @Tags C端-互动
// @Accept json
// @Produce json
// @Param id path int true "帖子ID"
// @Param body body CreateCommentReq true "content"
// @Success 200 {object} model.Comment
// @Router /api/v1/posts/{id}/comments [post]
func (h *CommentsHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	postID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if postID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req CreateCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content required"})
		return
	}
	comment := &model.Comment{PostID: postID, UserID: userID, Content: req.Content}
	if err := h.commentRepo.Create(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 非本人评论时，异步发 Kafka 消息，由消费者写库
	if post, _ := h.postRepo.GetByID(postID); post != nil && post.UserID != userID && post.UserID > 0 {
		_ = h.notifProd.Publish(context.Background(), post.UserID, "comment", comment.ID, model.JSONMap{"actor_id": userID, "post_id": postID})
	}
	u, _ := h.userRepo.GetByID(userID)
	comment.User = u
	c.JSON(http.StatusOK, comment)
}

type CreateCommentReq struct {
	Content string `json:"content" binding:"required"`
}
