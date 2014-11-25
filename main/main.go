package main

import (
	"github.com/mrjones/hanabi/hanabi"
	"github.com/mrjones/hanabi/strategies"

	"log"
)

func main() {
	state, err := hanabi.InitializeGame(
		[]hanabi.PlayerStrategy{
			&strategies.AlwaysPlayFirstCardStrategy{"Matt"},
			&strategies.AlwaysPlayFirstCardStrategy{"Cristina"},
		},
		[]hanabi.Observer{
			&hanabi.LoggingObserver{},
		})

	if err != nil {
		log.Fatal(err)
	}

	state.Play()
}
