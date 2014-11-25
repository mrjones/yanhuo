package main

import (
	"github.com/mrjones/hanabi/hanabi"
	"github.com/mrjones/hanabi/strategies"

	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello, world!")
	
	players := []hanabi.PlayerStrategy{
		&strategies.AlwaysPlayFirstCardStrategy{"Matt"},
		&strategies.AlwaysPlayFirstCardStrategy{"Cristina"}}

	observers := []hanabi.Observer{
		&hanabi.LoggingObserver{},
	}


	state, err := hanabi.InitializeGame(players, observers)
	
	if err != nil {
		log.Fatal(err)
	}

	for {
		state.TakeTurn()
	}
}
