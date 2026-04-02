package main

import (
	"fmt"

	"github.com/jondatkins/learn-pub-sub-starter/internal/gamelogic"
	"github.com/jondatkins/learn-pub-sub-starter/internal/pubsub"
	"github.com/jondatkins/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, channel *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(recOfWar gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")

		outcome, winner, loser := gs.HandleWar(recOfWar)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeYouWon:
			logMessage := fmt.Sprint("%s won a war against %s", loser, winner)
			PublishGameLog(channel, gs.GetUsername(), logMessage)
			return pubsub.Ack
		case gamelogic.WarOutcomeOpponentWon:
			logMessage := fmt.Sprint("%s won a war against %s", winner, loser)
			PublishGameLog(channel, gs.GetUsername(), logMessage)
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			logMessage := fmt.Sprint("A war between %s and % resulted in a draw", winner, loser)
			PublishGameLog(channel, gs.GetUsername(), logMessage)
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
				fmt.Printf("error:%s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
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
