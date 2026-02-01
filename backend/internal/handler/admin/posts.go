package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/internal/repository"
	"github.com/vibecoding/community/internal/service"
)

type AdminPostsHandler struct {
	post *service.PostService
	user *repository.UserRepo
}

func NewAdminPostsHandler(post *service.PostService, user *repository.UserRepo) *AdminPostsHandler {
	return &AdminPostsHandler{post: post, user: user}
}

// Get 管理后台帖子详情（含完整内容与媒体）
// @Summary 管理后台帖子详情
// @Tags B端-内容
// @Produce json
// @Param id path int true "帖子ID"
// @Success 200 {object} model.Post
// @Router /api/v1/admin/posts/{id} [get]
func (h *AdminPostsHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	p, err := h.post.GetByID(id, true, true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// List 管理后台帖子列表（含隐藏/标记）
// @Summary 管理后台帖子列表
// @Tags B端-内容
// @Produce json
// @Param offset query int false "0"
// @Param limit query int false "20"
// @Param status query string false "normal|hidden|flagged"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/posts [get]
func (h *AdminPostsHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	statusStr := c.DefaultQuery("status", "all") // 默认「全部」，隐藏/标记的帖子仍可见
	status := model.PostStatus(statusStr)
	list, total, err := h.post.List(offset, limit, status, "", 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// UpdateStatus 隐藏/标记帖子
// @Summary 更新帖子状态
// @Tags B端-内容
// @Accept json
// @Produce json
// @Param id path int true "帖子ID"
// @Param body body UpdatePostStatusReq true "status"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/posts/{id} [patch]
func (h *AdminPostsHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req UpdatePostStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := model.PostStatus(req.Status)
	if status != model.PostStatusNormal && status != model.PostStatusHidden && status != model.PostStatusFlagged {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	if err := h.post.UpdateStatus(id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Delete 删除帖子
// @Summary 删除帖子
// @Tags B端-内容
// @Param id path int true "帖子ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/posts/{id} [delete]
func (h *AdminPostsHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.post.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type UpdatePostStatusReq struct {
	Status string `json:"status"` // normal, hidden, flagged
}
