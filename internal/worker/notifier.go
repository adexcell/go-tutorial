package worker

import (
	"context"
	"encoding/json"

	"github.com/adexcell/go-tutorial/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func Start(ctx context.Context, msgs <-chan amqp.Delivery, logger zerolog.Logger) {
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("worker остановлен")
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}
			n := &domain.Notification{}
			if err := json.Unmarshal(msg.Body, &n); err != nil {
				logger.Error().Err(err).Msg("ошибка чтения сообщения")
				continue
			}
			logger.Info().Msgf("ОТПРАВКА УВЕДОМЛЕНИЯ: [User: %d] Message: %s", n.ID, n.Message)
			msg.Ack(false)
		}
	}
}
