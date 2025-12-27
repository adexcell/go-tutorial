package domain

import (
	"context"
	"time"
)

type Notification struct {
	ID      int64 `json:"id"`
	UserID  int64 `json:"user_id"`
	Message string `json:"message" binding:"required"`
	SendAt  time.Time `json:"send_at" binding:"required"`
}

type NotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
}

type NotificationSender interface {
	Publish(ctx context.Context, n *Notification) error
}

type NotificationService interface {
	Schedule(ctx context.Context, n *Notification) error
}