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
	fmt.Println("Starting Peril client...")

	connectString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatal("Client failed to connect to amqp")
	}
	defer connection.Close()
	fmt.Println("Connection Successful")
	username, err := ClientWelcome()
	if err != nil {
		fmt.Println("Error getting username")
		return
	}
	// fmt.Println(username)

	_, _, err = DeclareAndBind(connection, ExchangePerilDirect, PauseKey+"."+username, PauseKey, QueueTypeTransient)
	if err != nil {
		fmt.Println("Error in DeclareAndBind call", err)
		return
	}

	// Create a channel to receive OS signals
	done := make(chan os.Signal, 1)
	// Notify the channel on SIGINT (Ctrl+C)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	// fmt.Println("Waiting for Ctrl+C to exit...")
	// Block until a signal is received
	<-done
	fmt.Println("Received Ctrl+C, exiting gracefully...")
}
