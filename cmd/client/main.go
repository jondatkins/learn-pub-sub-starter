package main

import (
	"fmt"
	"log"
	"os"

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

	newGameState := NewGameState(username)
	// Add a REPL loop similar to what you did in the cmd/server application. Here's what each command should do:
	// The spawn command allows a player to add a new unit to the map under their control. Use the gamestate.CommandSpawn method and pass in the "words" from the GetInput command.
	// Possible unit types are: infantry, cavalry, artillery
	// Possible locations are: americas, europe, africa, asia, antarctica, australia
	// Example usage: spawn europe infantry
	// After spawning a unit, you should see its ID printed to the console.
	type Unit string

	// const (
	// 	Infantry  Unit = "infantry"
	// 	Cavalry   Unit = "cavalry"
	// 	Artillery Unit = "artillery"
	// )
	//
	// type Location string
	//
	// const (
	// 	Europe     Location = "europe"
	// 	Africa     Location = "africa"
	// 	Asia       Location = "asia"
	// 	Antarctica Location = "antarctica"
	// 	Australia  Location = "australia"
	// )
	//
	for {
		input := GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			newGameState.CommandSpawn(input)
		case "move":
			newGameState.CommandMove(input)
			if err == nil {
				fmt.Println("Move successful")
			}
		case "status":
			newGameState.CommandStatus()
		case "help":
			PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			PrintQuit()
			os.Exit(0)
			// Notify the channel on SIGINT (Ctrl+C)
		default:
			fmt.Println("Bad Input, expected 'spawn/move/status/help/spam/quit")
		}
	}
	// Create a channel to receive OS signals

	// // Notify the channel on SIGINT (Ctrl+C)
	// signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	// // Block until a signal is received
	// <-done
	// fmt.Println("Received Ctrl+C, exiting gracefully...")
}
