package c

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vibecoding/community/internal/middleware"
	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/internal/service"
)

type PostsHandler struct {
	post *service.PostService
}

func NewPostsHandler(post *service.PostService) *PostsHandler {
	return &PostsHandler{post: post}
}

// List 帖子列表，按时间倒序
// @Summary 帖子列表
// @Tags C端-帖子
// @Produce json
// @Param order query string false "desc"
// @Param offset query int false "0"
// @Param limit query int false "20"
// @Param keyword query string false "搜索"
// @Param tag query string false "话题名"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/posts [get]
func (h *PostsHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit > 100 {
		limit = 100
	}
	keyword := c.Query("keyword")
	tag := c.Query("tag")
	var topicID uint64
	if tag != "" {
		// 简化：这里可按 tag 名查 topic_id，暂时传 0 表示不按话题筛
		topicID = 0
	}
	list, total, err := h.post.List(offset, limit, model.PostStatusNormal, keyword, topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 附加点赞数、评论数（按帖子隔离：每条帖子的 like_count/comment_count 仅统计该帖子）
	type item struct {
		model.Post
		LikeCount    int64 `json:"like_count"`
		CommentCount int64 `json:"comment_count"`
	}
	out := make([]item, len(list))
	for i := range list {
		postID := list[i].ID
		lc, _ := h.post.LikeCount(postID)
		cc, _ := h.post.CommentCount(postID)
		out[i] = item{Post: list[i], LikeCount: lc, CommentCount: cc}
	}
	c.JSON(http.StatusOK, gin.H{"list": out, "total": total})
}

// Get 帖子详情（可选鉴权：带 token 时返回当前用户是否已点赞 liked）
// @Summary 帖子详情
// @Tags C端-帖子
// @Produce json
// @Param id path int true "帖子ID"
// @Success 200 {object} model.Post
// @Router /api/v1/posts/{id} [get]
func (h *PostsHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	p, err := h.post.GetByID(id, true, true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if p.Status != model.PostStatusNormal {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	likes, _ := h.post.LikeCount(id)
	comments, _ := h.post.CommentCount(id)
	liked := false
	if userID := middleware.GetUserID(c); userID > 0 {
		liked, _ = h.post.HasLiked(userID, id)
	}
	c.JSON(http.StatusOK, gin.H{
		"post":          p,
		"like_count":    likes,
		"comment_count": comments,
		"liked":         liked,
	})
}

// Create 发帖（需登录）
// @Summary 发帖
// @Tags C端-帖子
// @Accept json
// @Produce json
// @Param body body CreatePostReq true "content, type, media_urls, topic_names"
// @Success 200 {object} model.Post
// @Router /api/v1/posts [post]
func (h *PostsHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req CreatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Content == "" && len(req.MediaURLs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content or media required"})
		return
	}
	postType := model.PostTypeText
	if len(req.MediaURLs) > 0 {
		if req.Type == "video" {
			postType = model.PostTypeVideo
		} else {
			postType = model.PostTypeImage
		}
	}
	p, err := h.post.Create(userID, req.Content, postType, req.MediaURLs, req.TopicNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

type CreatePostReq struct {
	Content    string   `json:"content"`
	Type       string   `json:"type"` // text, image, video
	MediaURLs  []string `json:"media_urls"`
	TopicNames []string `json:"topic_names"`
}

// MyPosts 个人主页帖子列表
// @Summary 我的帖子
// @Tags C端-帖子
// @Produce json
// @Param offset query int false "0"
// @Param limit query int false "20"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/me/posts [get]
func (h *PostsHandler) MyPosts(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	list, total, err := h.post.ListByUserIDWithDetails(userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 按帖子隔离：每条帖子的 like_count/comment_count 仅统计该帖子
	type item struct {
		model.Post
		LikeCount    int64 `json:"like_count"`
		CommentCount int64 `json:"comment_count"`
	}
	out := make([]item, len(list))
	for i := range list {
		postID := list[i].ID
		lc, _ := h.post.LikeCount(postID)
		cc, _ := h.post.CommentCount(postID)
		out[i] = item{Post: list[i], LikeCount: lc, CommentCount: cc}
	}
	c.JSON(http.StatusOK, gin.H{"list": out, "total": total})
}
