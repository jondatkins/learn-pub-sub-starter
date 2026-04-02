package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	. "github.com/jondatkins/learn-pub-sub-starter/internal/gamelogic"
	"github.com/jondatkins/learn-pub-sub-starter/internal/pubsub"
	. "github.com/jondatkins/learn-pub-sub-starter/internal/pubsub"
	"github.com/jondatkins/learn-pub-sub-starter/internal/routing"
	. "github.com/jondatkins/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril client...")
	const connectString = "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer connection.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	publishCh, err := connection.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	username, err := ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}
	gameState := NewGameState(username)

	pubsub.SubscribeJSON(
		connection,
		ExchangePerilTopic,
		ArmyMovesPrefix+"."+gameState.GetUsername(),
		ArmyMovesPrefix+".*",
		SimpleQueueTransient,
		// NewHandler here
		handlerMove(gameState, publishCh),
	)
	if err != nil {
		log.Fatalf("could not subscribe to army moves: %v", err)
	}

	pubsub.SubscribeJSON(
		connection,
		ExchangePerilDirect,
		PauseKey+"."+username,
		PauseKey,
		SimpleQueueTransient,
		// NewHandler here
		handlerPause(gameState),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	pubsub.SubscribeJSON(
		connection,
		ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		SimpleQueueDurable,
		handlerWar(gameState, publishCh),
	)
	if err != nil {
		log.Fatalf("could not subscribe to war queue: %v", err)
	}

	for {
		input := GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		// The move command in the REPL should now publish a move.
		// Publish the move to the army_moves.username routing key, where username
		// is the name of the player.
		// Use the peril_topic exchange.
		// Log a message to the console stating that the move was published successfully.
		case "move":
			move, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// gameLogic := GameLog{
			// 	CurrentTime: time.Now(),
			// 	Message:     input[1],
			// 	Username:    username,
			// }
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+"."+username,
				move,
			)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// fmt.Println("Move published successfully: ", input)
			fmt.Printf("Moved %v units to %s\n", len(move.Units), move.ToLocation)

		case "spawn":
			gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}
}

// Create a reusable function to publish a GameLog struct:
// The topic exchange.
// The GameLogSlug.username routing key, where username is the name of the player who initiated the war, and GameLogSlug is a constant in the routing package.
// The GameLog struct should be serialized using the PublishGob function. Fill all the fields in.
// If a publishing fails, NackRequeue, otherwise Ack.
func PublishGameLog(ch *amqp.Channel, username, message string) error {
	return pubsub.PublishGob(
		ch,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+username,
		routing.GameLog{
			Username:    username,
			CurrentTime: time.Now(),
			Message:     message,
		},
	)
}
