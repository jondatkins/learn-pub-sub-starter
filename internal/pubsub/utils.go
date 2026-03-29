package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	QueueTypeDurable   SimpleQueueType = "durable"
	QueueTypeTransient SimpleQueueType = "transient"
)

// func (ch *Channel) PublishWithContext(_ context.Context, exchange, key string, mandatory, immediate bool, msg Publishing) error
func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	valBytes, err := json.Marshal(val)
	if err != nil {
		fmt.Println("Error marshalling value", err)
		return err
	}
	// Use the channel's .PublishWithContext method to publish
	// the message to the exchange with the routing key. Some configurations:
	// Set ctx to context.Background()
	// Set mandatory to false.
	// Set immediate to false.
	// In the amqp.Publishing struct, you only need to set two fields:
	// ContentType to "application/json".
	// Body to the JSON bytes.
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         valBytes,
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// This is not a mandatory delivery, so it will be dropped if there are no
	// queues bound to the logs exchange.
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
	if err != nil {
		// Since publish is asynchronous this can happen if the network connection
		// is reset or if the server has run out of resources.
		fmt.Println("Error publishing with context", err)
		return err
	}
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to create channel on connection")
	}
	dur := false
	if queueType == QueueTypeDurable {
		dur = true
	}
	autoDelete := false
	exclusive := false
	if queueType == QueueTypeTransient {
		autoDelete = true
		exclusive = true
	}
	q, _ := channel.QueueDeclare(
		queueName,  // queue name
		dur,        // durable
		autoDelete, // auto-delete
		exclusive,  // exclusive
		false,      // noWait
		nil,
	)
	log.Printf("Declared Classic Queue v2: %s", q.Name)
	fmt.Println("exchange is: ", exchange)
	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	return channel, q, nil
}
