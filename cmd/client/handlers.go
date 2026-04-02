package main

import (
	"fmt"

	"github.com/jondatkins/learn-pub-sub-starter/internal/gamelogic"
	"github.com/jondatkins/learn-pub-sub-starter/internal/pubsub"
	"github.com/jondatkins/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Create a new handler that consumes all the war messages that the "move" handler publishes, no matter the username in the routing key. It should:
// defer fmt.Print("> ") to ensure a new prompt is printed after the handler is done.
// Call the gamestate's HandleWar method with the message's body.
// If the outcome is gamelogic.WarOutcomeNotInvolved: NackRequeue the message so another client can try to consume it.
// If the outcome is gamelogic.WarOutcomeNoUnits: NackDiscard the message.
// If the outcome is gamelogic.WarOutcomeOpponentWon: Ack the message.
// If the outcome is gamelogic.WarOutcomeYouWon: Ack the message.
// If the outcome is gamelogic.WarOutcomeDraw: Ack the message.
// If it's anything else, print an error and NackDiscard the message.

func handlerWarMessages(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(recOfWar gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")

		outcome, _, _ := gs.HandleWar(recOfWar)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeYouWon:
			return pubsub.Ack
		case gamelogic.WarOutcomeOpponentWon:
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			return pubsub.Ack
		}
		fmt.Println("error: unknown war outcome")
		return pubsub.NackDiscard
	}
}

func handlerMove(gs *gamelogic.GameState, channel *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(armyMove gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(armyMove)
		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			fmt.Println("should be 'war.somName' ", routing.WarRecognitionsPrefix+"."+gs.GetUsername())

			err := pubsub.PublishJSON(
				channel,
				routing.ExchangePerilTopic,
				routing.WarRecognitionsPrefix+"."+gs.GetUsername(),
				gamelogic.RecognitionOfWar{
					Attacker: armyMove.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				fmt.Println(err)
			}
			return pubsub.NackRequeue
		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(playingState routing.PlayingState) pubsub.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(playingState)
		return pubsub.Ack
	}
}
