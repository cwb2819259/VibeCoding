package notification

import (
	"context"
	"fmt"
	"log"

	"github.com/vibecoding/community/internal/model"
	"github.com/vibecoding/community/pkg/kafka"
)

// Producer 异步发送通知到 Kafka，由消费者消费后写入 DB
type Producer struct {
	w *kafka.Writer
}

func NewProducer(w *kafka.Writer) *Producer {
	return &Producer{w: w}
}

// Publish 发送通知消息到 Kafka（异步，不阻塞接口）
func (p *Producer) Publish(ctx context.Context, userID uint64, typ string, relatedID uint64, payload model.JSONMap) error {
	log.Printf("[notif-producer] Publish: userID=%d type=%s relatedID=%d", userID, typ, relatedID)
	key := fmt.Sprintf("%d", userID)
	msg := model.NotificationMessage{
		UserID:    userID,
		Type:      typ,
		RelatedID: relatedID,
		Payload:   payload,
	}
	err := p.w.Send(ctx, key, &msg)
	if err != nil {
		log.Printf("[notif-producer] Send failed: %v", err)
	} else {
		log.Printf("[notif-producer] Send ok: userID=%d type=%s", userID, typ)
	}
	return err
}
