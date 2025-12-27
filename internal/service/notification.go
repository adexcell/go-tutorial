package service

import (
	"context"

	"github.com/adexcell/go-tutorial/internal/domain"
)

type NotificationService struct {
	sender domain.NotificationSender
}

func NewNotificationService(sender domain.NotificationSender) *NotificationService {
	return &NotificationService{sender: sender}
}

func (s *NotificationService) Schedule(ctx context.Context, n *domain.Notification) error {
	return s.sender.Publish(ctx, n)
}
