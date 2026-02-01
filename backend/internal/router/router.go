package router

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"github.com/vibecoding/community/internal/config"
	"github.com/vibecoding/community/internal/handler/admin"
	"github.com/vibecoding/community/internal/handler/c"
	"github.com/vibecoding/community/internal/middleware"
	"github.com/vibecoding/community/internal/notification"
	"github.com/vibecoding/community/internal/repository"
	"github.com/vibecoding/community/internal/service"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config, notifProd *notification.Producer) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	// doc 放在根路径，避免 /swagger/doc.json 与 /swagger/*any 在 Gin 中冲突
	r.GET("/doc.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		doc, err := swag.ReadDoc("swagger")
		if err != nil {
			c.String(500, `{"error":"swagger doc not found"}`)
			return
		}
		c.String(200, doc)
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/doc.json"), ginSwagger.InstanceName("swagger")))

	// 上传文件静态访问：GET /uploads/xxx 返回本地文件
	if absDir, err := filepath.Abs(cfg.Upload.Dir); err == nil {
		r.Static("/uploads", absDir)
	}

	userRepo := repository.NewUserRepo(db)
	adminRepo := repository.NewAdminRepo(db)
	postRepo := repository.NewPostRepo(db)
	likeRepo := repository.NewLikeRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	notificationRepo := repository.NewNotificationRepo(db)
	topicRepo := repository.NewTopicRepo(db)

	authUser := service.NewAuthUserService(userRepo, cfg.JWT.UserSecret, cfg.JWT.ExpireHours)
	authAdmin := service.NewAuthAdminService(adminRepo, cfg.JWT.AdminSecret, cfg.JWT.ExpireHours)
	postSvc := service.NewPostService(postRepo, userRepo, topicRepo, likeRepo, commentRepo)
	notifSvc := service.NewNotificationService(notificationRepo)

	// C 端
	authCH := c.NewAuthHandler(authUser)
	postsCH := c.NewPostsHandler(postSvc)
	uploadCH := c.NewUploadHandler(cfg.Upload.Dir, cfg.Upload.URLPrefix)
	likesCH := c.NewLikesHandler(likeRepo, postRepo, notifProd)
	commentsCH := c.NewCommentsHandler(commentRepo, postRepo, userRepo, notifProd)
	notifCH := c.NewNotificationsHandler(notifSvc)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", authCH.Login)
		v1.GET("/posts", postsCH.List)
		v1.GET("/posts/:id", middleware.OptionalAuthUser(cfg.JWT.UserSecret), postsCH.Get)
		v1.POST("/upload", uploadCH.Upload)

		needUser := v1.Group("")
		needUser.Use(middleware.AuthUser(cfg.JWT.UserSecret))
		{
			needUser.POST("/posts", postsCH.Create)
			needUser.GET("/users/me/posts", postsCH.MyPosts)
			needUser.POST("/posts/:id/like", likesCH.Like)
			needUser.DELETE("/posts/:id/like", likesCH.Unlike)
			needUser.POST("/posts/:id/comments", commentsCH.Create)
			needUser.GET("/notifications/unread-count", notifCH.UnreadCount)
			needUser.GET("/notifications", notifCH.List)
			needUser.PATCH("/notifications/:id/read", notifCH.MarkRead)
		}
		// 未登录也可看评论列表
		v1.GET("/posts/:id/comments", commentsCH.List)
	}

	// B 端
	adminAuth := admin.NewAuthHandler(authAdmin)
	adminPosts := admin.NewAdminPostsHandler(postSvc, userRepo)
	adminUsers := admin.NewAdminUsersHandler(userRepo)
	adminStats := admin.NewAdminStatsHandler(db)

	adminGroup := r.Group("/api/v1/admin")
	{
		adminGroup.POST("/auth/login", adminAuth.Login)
		needAdmin := adminGroup.Group("")
		needAdmin.Use(middleware.AuthAdmin(cfg.JWT.AdminSecret))
		{
			needAdmin.GET("/posts", adminPosts.List)
			needAdmin.GET("/posts/:id", adminPosts.Get)
			needAdmin.PATCH("/posts/:id", adminPosts.UpdateStatus)
			needAdmin.DELETE("/posts/:id", adminPosts.Delete)
			needAdmin.GET("/users", adminUsers.List)
			needAdmin.GET("/stats", adminStats.Stats)
		}
	}

	return r
}
