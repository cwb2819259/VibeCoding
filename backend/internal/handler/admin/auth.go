package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vibecoding/community/internal/service"
)

type AuthHandler struct {
	auth *service.AuthAdminService
}

func NewAuthHandler(auth *service.AuthAdminService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type AdminLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResp struct {
	Token string      `json:"token"`
	Admin interface{} `json:"admin"`
}

// Login B端管理员登录
// @Summary 管理员登录
// @Tags B端-认证
// @Accept json
// @Produce json
// @Param body body AdminLoginReq true "username, password"
// @Success 200 {object} AdminLoginResp
// @Router /api/v1/admin/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, admin, err := h.auth.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, AdminLoginResp{Token: token, Admin: admin})
}
