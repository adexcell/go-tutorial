package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/adexcell/go-tutorial/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationQueue struct {
	conn *amqp.Connection
}

func NewNotificationQueue(conn *amqp.Connection) *NotificationQueue {
	return &NotificationQueue{conn: conn}
}

func (q *NotificationQueue) Publish(ctx context.Context, n *domain.Notification) error {
	ch, err := q.conn.Channel()
	if err != nil {
		return fmt.Errorf("не удалось открыть канал: %w", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"delayed_exchange",  // имя
		"x-delayed-message", // ТИП
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		amqp.Table{
			"x-delayed-type": "direct", // какой тип будет у "внутреннего" обмена
		},
	)
	if err != nil {
		return fmt.Errorf("не удалось объявить Exchange%w", err)
	}

	delay := time.Until(n.SendAt).Milliseconds()
	if delay < 0 { 
		delay = 0 
	}

	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("не удалось подготовить данные для публикации в брокере: %w", err)
	}

	if err := ch.PublishWithContext(ctx, "delayed_exchange", "notification_key", false, false, amqp.Publishing{
		Headers: amqp.Table{ "x-delay": delay},
		ContentType: "application/json",
		Body: body,
	}); err != nil {
		return fmt.Errorf("не удалось сделать публикацию в брокер: %w", err)
	}

	return nil
}

func (q *NotificationQueue) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	ch, err := q.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть канал: %w", err)
	}

	queue, err := ch.QueueDeclare(
		"notifications_queue", // name
        true,                  // durable
        false,                 // delete when unused
        false,                 // exclusive
        false,                 // no-wait
        nil,                   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось объявить Queue%w", err)
	}

	err = ch.QueueBind(
        queue.Name,           // queue name
        "notification_key",   // routing key (тот же, что в Publish)
        "delayed_exchange",   // exchange
        false,
        nil,
	)

	return ch.Consume(
        queue.Name, // queue
        "",         // consumer
        false,      // auto-ack (ставим false, затем вручную вызываем в воркере)
        false,      // exclusive
        false,      // no-local
        false,      // no-wait
        nil,        // args
    )
}