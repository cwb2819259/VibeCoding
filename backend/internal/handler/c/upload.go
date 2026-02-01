package c

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vibecoding/community/internal/middleware"
)

type UploadHandler struct {
	uploadDir    string
	uploadPrefix string
}

func NewUploadHandler(uploadDir, uploadPrefix string) *UploadHandler {
	return &UploadHandler{uploadDir: uploadDir, uploadPrefix: uploadPrefix}
}

// Upload 上传图片/视频，返回 URL
// @Summary 上传文件
// @Tags C端-上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} map[string]string
// @Router /api/v1/upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	_ = middleware.GetUserID(c) // 可选：要求登录后才允许上传
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
		".mp4": true, ".webm": true,
	}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
		return
	}
	name := uuid.New().String() + ext
	dst := filepath.Join(h.uploadDir, name)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	url := h.uploadPrefix + "/" + name
	c.JSON(http.StatusOK, gin.H{"url": url})
}
