package main

import (
	"github.com/mrjones/hanabi/hanabi"

	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello, world!")
	
	deck := hanabi.Shuffle(hanabi.CreateDeck())
	players := []hanabi.PlayerStrategy{
		&hanabi.AlwaysPlayFirstCardStrategy{"Matt"},
		&hanabi.AlwaysPlayFirstCardStrategy{"Cristina"}}

	state, err := hanabi.InitializeGame(deck, players)
	
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
