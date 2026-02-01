package c

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vibecoding/community/internal/service"
)

type AuthHandler struct {
	auth *service.AuthUserService
}

func NewAuthHandler(auth *service.AuthUserService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type LoginReq struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type LoginResp struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// Login C端登录：手机号+验证码，mock 固定 123456
// @Summary C端登录
// @Tags C端-认证
// @Accept json
// @Produce json
// @Param body body LoginReq true "phone, code"
// @Success 200 {object} LoginResp
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, user, err := h.auth.Login(req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, LoginResp{Token: token, User: user})
}
