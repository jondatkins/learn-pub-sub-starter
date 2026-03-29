package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"

	. "github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	. "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	. "github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatal("Failed to connect to amqp")
	}
	defer connection.Close()
	fmt.Println("Connection Successful")

	// Create a channel to receive OS signals
	done := make(chan os.Signal, 1)

	// Notify the channel on SIGINT (Ctrl+C)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	PrintServerHelp()
	// At the beginning of the loop, use the GetInput function in internal/gamelogic to wait
	// for a slice of input "words" from the user. If the slice is empty, continue to the next iteration of the loop.
	// Check the first word:
	// If it's "pause", log to the console that you're sending a pause message, and publish the pause message as you were doing before.
	// If it's "resume", log to the console that you're sending a resume message, and publish the resume message as you were doing before.
	// The only difference is that the IsPaused field should be set to false.
	// If it's "quit", log to the console that you're exiting, and break out of the loop.
	// If it's anything else, log to the console that you don't understand the command.
	for {
		input := GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			fmt.Println("Sending Pause Message")
			publishMessage(connection, PauseKey, true)
		case "resume":
			fmt.Println("Sending Resume Message")
			publishMessage(connection, ResumeKey, false)
		case "quit":
			fmt.Println("Exiting")
			break
		default:
			fmt.Println("Bad Command")
		}
	}

	// Block until a signal is received
	// <-done
	// fmt.Println("Received Ctrl+C, exiting gracefully...")
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
