package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Writer struct {
	w *kafka.Writer
}

func NewWriter(brokers []string, topic string) *Writer {
	return &Writer{
		w: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (w *Writer) Send(ctx context.Context, key string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return w.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: b,
	})
}

func (w *Writer) Close() error {
	return w.w.Close()
}
