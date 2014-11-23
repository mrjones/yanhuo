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

	state, err := hanabi.InitializeGame(players)
	
	if err != nil {
		log.Fatal(err)
	}

	/*
	for i := hanabi.PlayerIndex(0); int(i) < len(players); i++ {
		fmt.Printf("PLAYER %d\n", i)
		hanabi.DisplayDeck(state.playerStates[i].cards)
	}
*/

	for {
		log.Println("---")
		state.TakeTurn()
	}
}
