package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Key   []byte
	Value []byte
	Topic string
}

type HandlerFunc func(Message) error

type ConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type Consumer struct {
	r *kafka.Reader
}

func NewConsumer(cfg ConsumerConfig) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  cfg.GroupID,
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})
	return &Consumer{r: r}
}

func (c *Consumer) Start(ctx context.Context, h HandlerFunc) error {
	for {
		m, err := c.r.ReadMessage(ctx)
		if err != nil {
			return err
		}

		msg := Message{
			Key:   m.Key,
			Value: m.Value,
			Topic: m.Topic,
		}

		if err := h(msg); err != nil {
			log.Printf("[kafka] handler error: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.r.Close()
}
