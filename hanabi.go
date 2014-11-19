package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Color int8
type Value int8
type PlayerIndex int8
type HandIndex int8

const (
	WHITE Color = iota
	RED
	BLUE
	YELLOW
	GREEN
)

var ALL_COLORS = []Color{WHITE, RED, BLUE, YELLOW, GREEN}

type colorInfo struct {
	fullName string
	shortName string
}

var COLOR_INFOS = map[Color]colorInfo {
	WHITE: colorInfo{fullName: "WHITE", shortName: "W"},
	RED: colorInfo{fullName: "RED", shortName: "R"},
	BLUE: colorInfo{fullName: "BLUE", shortName: "B"},
	YELLOW: colorInfo{fullName: "YELLOW", shortName: "Y"},
	GREEN: colorInfo{fullName: "GREEN", shortName: "G"},
}


var VALUE_COUNTS = map[Value]int{
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

type GiveInformationAction struct {
	PlayerIndex PlayerIndex

	// Exactly one of color or value must be non-null
	Color *Color
	Value *Value
}

type DiscardAction struct {
	// The index of the card to discard
	Index HandIndex
}

type PlayAction struct {
	// The index of the card to discard
	Index HandIndex
}

type Action struct {
	// Exactly one one must be non-null
	GiveInformation *GiveInformationAction
	Discard *DiscardAction
	Play *PlayAction
}

func (a Action) DebugString() string {
	switch {
	case a.GiveInformation != nil:
		return fmt.Sprintf("Gave information to player: %d", a.GiveInformation.PlayerIndex)
	case a.Discard != nil:
		return fmt.Sprintf("Discarded card %d", a.Discard.Index)
	case a.Play != nil:
		return fmt.Sprintf("Played card %d", a.Play.Index)
	default:
		return "INVALID ACTION"
	}
}


type PlayerLogic interface {
	// map player -> []card
	// Act
	Act(otherPlayersCards map[PlayerIndex][]Card, myNumCards int, blueTokens int, redTokens int) Action

	ObserveAction(actor PlayerIndex, action Action)
}

type SimplePlayerLogic struct {
	Name string
}

func (p *SimplePlayerLogic) Act(otherPlayersCards map[PlayerIndex][]Card, myNumCards int, blueTokens int, redTokens int) Action {
	return Action{Play: &PlayAction{Index: 0}}
}

func (p *SimplePlayerLogic) ObserveAction(actor PlayerIndex, action Action) {
	log.Printf("%s observed %s (by %d)\n", p.Name, action, actor)
}

func main() {
	fmt.Println("Hello, world!")
	
	deck := shuffle(createDeck())
	players := []PlayerLogic{
		&SimplePlayerLogic{"Matt"}, &SimplePlayerLogic{"Cristina"}}

	state, err := initializeGame(deck, players)
	
	if err != nil {
		log.Fatal(err)
	}

	for i := PlayerIndex(0); int(i) < len(players); i++ {
		fmt.Printf("PLAYER %d\n", i)
		displayDeck(state.playerStates[i].cards)
	}
}

type Card struct {
	Value Value
	Color Color
}

type playerState struct {
	cards []Card
	logic PlayerLogic
}

type gameState struct {
	drawPile []Card
	playerStates []playerState

	redTokens int  // bad plays
	blueTokens int  // available information

	
}

func initializeGame(deck []Card, players []PlayerLogic) (*gameState, error) {
	numPlayers := len(players)

	cardsPerPlayer, ok := INITIAL_CARDS[numPlayers]
	if !ok {
		return nil, fmt.Errorf("Invalid number of players: %d", numPlayers)
	}

	state := &gameState{
		playerStates: make([]playerState, numPlayers),
		drawPile: []Card{},
	}

	for p := PlayerIndex(0); int(p) < numPlayers; p++ {
		state.playerStates[p] = playerState{
			cards: []Card{},
			logic: players[p],
		}
	}

	drawCount := 0
	for c := 0; c < cardsPerPlayer; c++ {
		for p := PlayerIndex(0); int(p) < numPlayers; p++ {
			// TODO(mrjones): bounds check
			state.playerStates[p].cards = append(state.playerStates[p].cards, deck[drawCount])
			drawCount++
		}
	}

	state.drawPile = deck[drawCount:]
	return state, nil
}

func displayDeck(deck []Card) {
	for _, card := range(deck) {
		fmt.Printf("color: %s, value: %d\n", COLOR_INFOS[card.Color].fullName, card.Value)
	}
}

func shuffle(in []Card) []Card {
	rand.Seed(time.Now().UTC().UnixNano())

	out := []Card{}
	for _, i := range rand.Perm(len(in)) {
		out = append(out, in[i])
	}

	return out
}

func createDeck() []Card {
	cards := []Card{}
	for color, _ := range(COLOR_INFOS) {
		for value, count := range(VALUE_COUNTS) {
			for i := 0; i < count; i++ {
				cards = append(cards, Card{Value: value, Color: color})
			}
		}
	}
	
	return cards
}
