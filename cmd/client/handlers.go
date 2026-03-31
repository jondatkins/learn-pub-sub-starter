package main

import (
	"fmt"

	"github.com/jondatkins/learn-pub-sub-starter/internal/gamelogic"
	"github.com/jondatkins/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(playingState routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(playingState)
	}
}
