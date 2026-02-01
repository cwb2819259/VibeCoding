package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vibecoding/community/internal/repository"
)

type AdminUsersHandler struct {
	user *repository.UserRepo
}

func NewAdminUsersHandler(user *repository.UserRepo) *AdminUsersHandler {
	return &AdminUsersHandler{user: user}
}

// List C端用户列表（users 表）
// @Summary 用户列表
// @Tags B端-用户
// @Produce json
// @Param offset query int false "0"
// @Param limit query int false "20"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/users [get]
func (h *AdminUsersHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	list, total, err := h.user.List(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}
