package consumers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hassiimykyta/life-rpg/pkg/kafka"
	usereventsv1 "github.com/hassiimykyta/life-rpg/services/events/user/v1"
)

type WelcomeSender interface {
	SendWelcome(to, username string) error
}

type UserRegistered struct {
	c       *kafka.Consumer
	handler WelcomeSender
}

func NewUserRegistered(brokers []string, groupID string, h WelcomeSender) *UserRegistered {
	return &UserRegistered{
		c: kafka.NewConsumer(kafka.ConsumerConfig{
			Brokers: brokers,
			Topic:   "user.registered",
			GroupID: groupID,
		}),
		handler: h,
	}
}

func (u *UserRegistered) Start(ctx context.Context) error {
	return u.c.Start(ctx, func(m kafka.Message) error {
		var evt usereventsv1.UserRegistered
		if err := json.Unmarshal(m.Value, &evt); err != nil {
			log.Printf("[user.registered] bad payload: %v", err)
			return nil
		}
		if err := u.handler.SendWelcome(evt.Email, evt.Username); err != nil {
			log.Printf("[user.registered] send welcome failed: %v", err)
		}
		return nil
	})
}

func (u *UserRegistered) Close() error { return u.c.Close() }
