package notification

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/internal/service"
)

// Consumer 消费 Kafka 通知消息并写入 DB
type Consumer struct {
	reader *kafka.Reader
	notif  *service.NotificationService
}

func NewConsumer(brokers []string, topic, groupID string, notif *service.NotificationService) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		notif: notif,
	}
}

// Run 阻塞消费，应在 goroutine 中调用；ctx 取消时退出
func (c *Consumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("[notif-consumer] ReadMessage error: %v", err)
			continue
		}
		log.Printf("[notif-consumer] message received, valueLen=%d", len(msg.Value))
		var m model.NotificationMessage
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			log.Printf("[notif-consumer] unmarshal error: %v", err)
			continue
		}
		log.Printf("[notif-consumer] unmarshaled: userID=%d type=%s relatedID=%d", m.UserID, m.Type, m.RelatedID)
		if m.UserID == 0 {
			log.Printf("[notif-consumer] skip user_id=0")
			continue
		}
		if err := c.notif.Create(m.UserID, m.Type, m.RelatedID, m.Payload); err != nil {
			log.Printf("[notif-consumer] Create error: %v", err)
			continue
		}
		log.Printf("[notif-consumer] created notification ok: user_id=%d type=%s related_id=%d", m.UserID, m.Type, m.RelatedID)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
