package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vibecoding/community/internal/model"
	"gorm.io/gorm"
)

type AdminStatsHandler struct {
	db *gorm.DB
}

func NewAdminStatsHandler(db *gorm.DB) *AdminStatsHandler {
	return &AdminStatsHandler{db: db}
}

// Stats 数据统计：发帖数、用户数等
// @Summary 数据统计
// @Tags B端-统计
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/stats [get]
func (h *AdminStatsHandler) Stats(c *gin.Context) {
	var userCount, postCount int64
	_ = h.db.Model(&model.User{}).Count(&userCount).Error
	_ = h.db.Model(&model.Post{}).Count(&postCount).Error
	c.JSON(http.StatusOK, gin.H{
		"user_count": userCount,
		"post_count": postCount,
	})
}
