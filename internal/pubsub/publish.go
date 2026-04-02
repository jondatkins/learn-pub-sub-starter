package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"time"

	"github.com/jondatkins/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	dat, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(
		context.Background(),
		exchange, key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        dat,
		})
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return err
	}
	return ch.PublishWithContext(
		context.Background(),
		exchange, key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        buf.Bytes(),
		})
}

// Create a reusable function to publish a GameLog struct:
// The topic exchange.
// The GameLogSlug.username routing key, where username is the name of the player who initiated the war, and GameLogSlug is a constant in the routing package.
// The GameLog struct should be serialized using the PublishGob function. Fill all the fields in.
// If a publishing fails, NackRequeue, otherwise Ack.
func PublishGameLog(ch *amqp.Channel, username, message string) (Acktype, error) {
	exchange := routing.ExchangePerilTopic
	key := routing.GameLogSlug + "." + username
	gameLog := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     message,
		Username:    username,
	}
	if err := PublishGob(ch, exchange, key, gameLog); err != nil {
		return NackRequeue, err
	}
	return Ack, nil
}
