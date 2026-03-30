package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/jondatkins/learn-pub-sub-starter/internal/gamelogic"
	. "github.com/jondatkins/learn-pub-sub-starter/internal/gamelogic"
	"github.com/jondatkins/learn-pub-sub-starter/internal/pubsub"
	. "github.com/jondatkins/learn-pub-sub-starter/internal/pubsub"
	"github.com/jondatkins/learn-pub-sub-starter/internal/routing"
	. "github.com/jondatkins/learn-pub-sub-starter/internal/routing"
)

func main() {
	const connectString = "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v")
	}
	defer connection.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	publishCh, err := connection.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		SimpleQueueDurable,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			fmt.Println("Publishing paused game state")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resumes game state")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "quit":
			log.Println("goodbye")
			return
		default:
			fmt.Println("unknown command")
		}
	}
}

func publishMessage(connection *amqp.Connection, key string, isPaused bool) error {
	newChannel, err := connection.Channel()
	if err != nil {
		log.Fatal("Failed to create channel on connection")
	}
	err = newChannel.ExchangeDeclare(ExchangePerilDirect, "direct", true, false, false, false, nil)
	playingState := PlayingState{IsPaused: isPaused}
	err = PublishJSON(newChannel, ExchangePerilDirect, key, playingState)
	if err != nil {
		fmt.Println("Error publishing JSON", err)
		return err
	}
	return nil
}
