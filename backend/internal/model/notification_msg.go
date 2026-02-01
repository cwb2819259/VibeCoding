package model

// NotificationMessage Kafka 通知消息体，由生产者发送、消费者消费后写入 DB
type NotificationMessage struct {
	UserID    uint64  `json:"user_id"`
	Type      string  `json:"type"`
	RelatedID uint64  `json:"related_id"`
	Payload   JSONMap `json:"payload,omitempty"`
}
