package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	. "github.com/jondatkins/learn-pub-sub-starter/internal/gamelogic"
	. "github.com/jondatkins/learn-pub-sub-starter/internal/pubsub"
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

	username, err := ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}

	_, queue, err := DeclareAndBind(
		connection,
		ExchangePerilDirect,
		PauseKey+"."+username,
		PauseKey,
		SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and not bound!\n", queue.Name)

	gameState := NewGameState(username)

	for {
		input := GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "move":
			gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
				continue
			}

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
