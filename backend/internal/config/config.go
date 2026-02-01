package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   Server
	MySQL    MySQL
	JWT      JWT
	Upload   Upload
	Kafka    Kafka
	Log      Log
}

type Log struct {
	Dir string // 日志目录，如 ./log 或 /app/log
}

type Server struct {
	Port string
}

type MySQL struct {
	DSN string
}

type JWT struct {
	UserSecret  string // C 端 JWT 密钥
	AdminSecret string // B 端 JWT 密钥
	ExpireHours int
}

type Upload struct {
	Dir    string // 本地存储目录，如 /data/upload
	URLPrefix string // 返回给前端的 URL 前缀，如 /uploads
}

type Kafka struct {
	Brokers []string
	TopicNotifications string
}

func Load() *Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		mysqlDSN = "root:root@tcp(localhost:3306)/community?charset=utf8mb4&parseTime=True"
	}
	userSecret := os.Getenv("JWT_USER_SECRET")
	if userSecret == "" {
		userSecret = "user-jwt-secret-change-in-prod"
	}
	adminSecret := os.Getenv("JWT_ADMIN_SECRET")
	if adminSecret == "" {
		adminSecret = "admin-jwt-secret-change-in-prod"
	}
	expireHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))
	if expireHours <= 0 {
		expireHours = 72
	}
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./data/upload"
	}
	uploadPrefix := os.Getenv("UPLOAD_URL_PREFIX")
	if uploadPrefix == "" {
		uploadPrefix = "/uploads"
	}
	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersEnv == "" {
		kafkaBrokersEnv = "localhost:9092"
	}
	brokers := strings.Split(kafkaBrokersEnv, ",")
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}
	topicNotif := os.Getenv("KAFKA_TOPIC_NOTIFICATIONS")
	if topicNotif == "" {
		topicNotif = "community.notifications"
	}
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "log"
	}

	return &Config{
		Server: Server{Port: port},
		MySQL:  MySQL{DSN: mysqlDSN},
		JWT: JWT{
			UserSecret:  userSecret,
			AdminSecret: adminSecret,
			ExpireHours: expireHours,
		},
		Upload: Upload{Dir: uploadDir, URLPrefix: uploadPrefix},
		Kafka: Kafka{
			Brokers:             brokers,
			TopicNotifications: topicNotif,
		},
		Log: Log{Dir: logDir},
	}
}
