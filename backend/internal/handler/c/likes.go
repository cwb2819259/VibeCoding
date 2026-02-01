package c

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vibecoding/community/internal/middleware"
	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/internal/notification"
	"github.com/vibecoding/community/internal/repository"
)

type LikesHandler struct {
	likeRepo  *repository.LikeRepo
	postRepo  *repository.PostRepo
	notifProd *notification.Producer
}

func NewLikesHandler(likeRepo *repository.LikeRepo, postRepo *repository.PostRepo, notifProd *notification.Producer) *LikesHandler {
	return &LikesHandler{likeRepo: likeRepo, postRepo: postRepo, notifProd: notifProd}
}

// Like 点赞帖子
// @Summary 点赞
// @Tags C端-互动
// @Produce json
// @Param id path int true "帖子ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/posts/{id}/like [post]
func (h *LikesHandler) Like(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	exists, _ := h.likeRepo.Exists(userID, id, model.TargetTypePost)
	if exists {
		count, _ := h.likeRepo.CountByTarget(id, model.TargetTypePost)
		c.JSON(http.StatusOK, gin.H{"liked": true, "like_count": count})
		return
	}
	l := &model.Like{UserID: userID, TargetType: model.TargetTypePost, TargetID: id}
	if err := h.likeRepo.Create(l); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 非本人点赞时，异步发 Kafka 消息，由消费者写库
	post, _ := h.postRepo.GetByID(id)
	if post == nil {
		log.Printf("[like] post not found, postID=%d, skip notification", id)
	} else if post.UserID == userID {
		log.Printf("[like] self like, postID=%d owner=%d actor=%d, skip notification", id, post.UserID, userID)
	} else if post.UserID == 0 {
		log.Printf("[like] post owner is 0, postID=%d, skip notification", id)
	} else {
		log.Printf("[like] sending notification: postID=%d owner=%d actor=%d", id, post.UserID, userID)
		if err := h.notifProd.Publish(context.Background(), post.UserID, "like", id, model.JSONMap{"actor_id": userID}); err != nil {
			log.Printf("[like] publish notification failed: %v", err)
		} else {
			log.Printf("[like] publish notification ok: owner=%d", post.UserID)
		}
	}
	count, _ := h.likeRepo.CountByTarget(id, model.TargetTypePost)
	c.JSON(http.StatusOK, gin.H{"liked": true, "like_count": count})
}

// Unlike 取消点赞
// @Summary 取消点赞
// @Tags C端-互动
// @Param id path int true "帖子ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/posts/{id}/like [delete]
func (h *LikesHandler) Unlike(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	_ = h.likeRepo.Delete(userID, id, model.TargetTypePost)
	count, _ := h.likeRepo.CountByTarget(id, model.TargetTypePost)
	c.JSON(http.StatusOK, gin.H{"liked": false, "like_count": count})
}
