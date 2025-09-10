package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

type ProducerConfig struct {
	Brokers      []string
	Topic        string
	RequiredAcks kafka.RequiredAcks
	Balancer     kafka.Balancer
}

func NewProducer(cfg ProducerConfig) *Producer {
	acks := cfg.RequiredAcks
	if acks == 0 {
		acks = kafka.RequireAll
	}
	bal := cfg.Balancer
	if bal == nil {
		bal = &kafka.LeastBytes{}
	}

	return &Producer{
		w: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        cfg.Topic,
			RequiredAcks: acks,
			Balancer:     bal,
		},
	}
}

func (p *Producer) Send(ctx context.Context, key, value []byte, headers ...kafka.Header) error {
	return p.w.WriteMessages(ctx, kafka.Message{
		Key:     key,
		Value:   value,
		Time:    time.Now(),
		Headers: headers,
	})
}

func (p *Producer) Close() error {
	return p.w.Close()
}
