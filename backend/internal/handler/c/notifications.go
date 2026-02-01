package c

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vibecoding/community/internal/middleware"
	"github.com/vibecoding/community/internal/service"
)

type NotificationsHandler struct {
	notif *service.NotificationService
}

func NewNotificationsHandler(notif *service.NotificationService) *NotificationsHandler {
	return &NotificationsHandler{notif: notif}
}

// UnreadCount 未读数量（供前端轮询）
// @Summary 未读数量
// @Tags C端-通知
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/notifications/unread-count [get]
func (h *NotificationsHandler) UnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	count, err := h.notif.CountUnread(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// List 消息通知列表
// @Summary 通知列表
// @Tags C端-通知
// @Produce json
// @Param offset query int false "0"
// @Param limit query int false "20"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/notifications [get]
func (h *NotificationsHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	list, total, err := h.notif.ListByUserID(userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"list": list, "total": total})
}

// MarkRead 标记已读
// @Summary 标记已读
// @Tags C端-通知
// @Param id path int true "通知ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/notifications/{id}/read [patch]
func (h *NotificationsHandler) MarkRead(c *gin.Context) {
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
	if err := h.notif.MarkRead(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
