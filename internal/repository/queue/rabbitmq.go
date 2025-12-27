package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func New(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к rabbitmq: %w", err)
	}
	return conn, nil
}
