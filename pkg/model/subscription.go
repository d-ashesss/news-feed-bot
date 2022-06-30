package model

import "context"

type Subscription struct {
	Category   Category
	Subscribed bool
	Unread     int
}

type SubscriptionModel interface {
	CreateSubscriber(ctx context.Context, s *Subscriber) (string, error)
	GetSubscriber(ctx context.Context, id string) (*Subscriber, error)
	Subscribe(ctx context.Context, s *Subscriber, c Category) error
	Unsubscribe(ctx context.Context, s *Subscriber, c Category) error
	GetSubscriptionStatus(ctx context.Context, s *Subscriber) ([]Subscription, error)
}
