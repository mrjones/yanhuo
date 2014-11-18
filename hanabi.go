package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type color int8
type value int8

const (
	WHITE color = iota
	RED
	BLUE
	YELLOW
	GREEN

)

var ALL_COLORS = []color{WHITE, RED, BLUE, YELLOW, GREEN}

type colorInfo struct {
	fullName string
	shortName string
}

var COLOR_INFOS = map[color]colorInfo {
	WHITE: colorInfo{fullName: "WHITE", shortName: "W"},
	RED: colorInfo{fullName: "RED", shortName: "R"},
	BLUE: colorInfo{fullName: "BLUE", shortName: "B"},
	YELLOW: colorInfo{fullName: "YELLOW", shortName: "Y"},
	GREEN: colorInfo{fullName: "GREEN", shortName: "G"},
}


var VALUE_COUNTS = map[value]int{
	1: 3,
	2: 2,
	3: 2,
	4: 2,
	5: 1,
}

// map from number of players to number of initial cards
var INITIAL_CARDS = map[int]int {
	2: 5,
	3: 5,
	4: 4,
	5: 4,
}

func main() {
	fmt.Println("Hello, world!")
	
	deck := shuffle(createDeck())
	state, err := initialDraw(deck, 4)
	
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 4; i++ {
		fmt.Printf("PLAYER %d\n", i)
		displayDeck(state.heldCards[i])
	}

}



type card struct {
	value value
	color color
}

type gameState struct {
	drawPile []card
	heldCards map[int][]card
}

func initialDraw(deck []card, numPlayers int) (*gameState, error) {
	cardsPerPlayer, ok := INITIAL_CARDS[numPlayers]
	if !ok {
		return nil, fmt.Errorf("Invalid number of players: %d", numPlayers)
	}

	state := &gameState{heldCards: make(map[int][]card), drawPile: []card{}}

	for p := 0; p < numPlayers; p++ {
		state.heldCards[p] = []card{}
	}

	drawCount := 0
	for c := 0; c < cardsPerPlayer; c++ {
		for p := 0; p < numPlayers; p++ {
			// TODO(mrjones): bounds check
			state.heldCards[p] = append(state.heldCards[p], deck[drawCount])
			drawCount++
		}
	}

	state.drawPile = deck[drawCount:]
	return state, nil
}

func displayDeck(deck []card) {
	for _, card := range(deck) {
		fmt.Printf("color: %s, value: %d\n", COLOR_INFOS[card.color].fullName, card.value)
	}
}

func shuffle(in []card) []card {
	rand.Seed(time.Now().UTC().UnixNano())

	out := []card{}
	for _, i := range rand.Perm(len(in)) {
		out = append(out, in[i])
	}

	return out
}

func createDeck() []card {
	cards := []card{}
	for color, _ := range(COLOR_INFOS) {
		for value, count := range(VALUE_COUNTS) {
			for i := 0; i < count; i++ {
				cards = append(cards, card{value: value, color: color})
			}
		}
	}
	
	return cards
}
