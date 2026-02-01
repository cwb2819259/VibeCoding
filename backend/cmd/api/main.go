package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/vibecoding/community/docs"
	"github.com/vibecoding/community/internal/config"
	"github.com/vibecoding/community/internal/notification"
	"github.com/vibecoding/community/internal/repository"
	"github.com/vibecoding/community/internal/router"
	"github.com/vibecoding/community/internal/service"
	"github.com/vibecoding/community/pkg/kafka"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title           VibeCoding 社区 API
// @version         1.0
// @description     VibeCoding社区 C 端（帖子、评论、点赞、通知、上传）与 B 端（内容管理、用户管理、数据统计）接口文档
// @termsOfService  https://github.com/vibecoding/community

// @contact.name   API Support
// @contact.url    https://github.com/vibecoding/community

// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT

// @host            localhost:8080
// @BasePath        /
// @schemes         http https
func main() {
	cfg := config.Load()
	// 日志同时输出到 stderr 和 log 目录下的 app.log
	if cfg.Log.Dir != "" {
		if err := os.MkdirAll(cfg.Log.Dir, 0755); err != nil {
			log.Println("log dir mkdir:", err)
		} else {
			logPath := filepath.Join(cfg.Log.Dir, "app.log")
			f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Println("log file open:", err)
			} else {
				log.SetOutput(io.MultiWriter(os.Stderr, f))
			}
		}
	}
	if err := os.MkdirAll(cfg.Upload.Dir, 0755); err != nil {
		log.Println("upload dir mkdir:", err)
	}
	// GORM SQL 日志使用与 log 相同的输出（终端 + log/app.log）
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.Fatal("mysql open:", err)
	}
	if err := service.SeedAdminIfNotExists(repository.NewAdminRepo(db)); err != nil {
		log.Fatal("seed admin:", err)
	}

	// Kafka 通知：生产者 + 消费者
	notifWriter := kafka.NewWriter(cfg.Kafka.Brokers, cfg.Kafka.TopicNotifications)
	defer notifWriter.Close()
	notifProd := notification.NewProducer(notifWriter)
	notifSvc := service.NewNotificationService(repository.NewNotificationRepo(db))
	notifConsumer := notification.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.TopicNotifications, "community-notif-consumer", notifSvc)
	defer notifConsumer.Close()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go notifConsumer.Run(ctx)

	r := router.Setup(db, cfg, notifProd)
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Println("listen", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
	_ = os.Stdin
}
