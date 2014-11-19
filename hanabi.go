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
	MAX_BLUE_TOKENS = 8

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
	// The player information is being given about
	PlayerIndex PlayerIndex

	// The cards matching the information being given
	Cards []HandIndex

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


type PlayerStrategy interface {
	Act(
		otherPlayersCards map[PlayerIndex][]Card,
		myNumCards int,
		blueTokens int,
		redTokens int) Action

	ObserveAction(actor PlayerIndex, action Action)
}

type AlwaysPlayFirstCardStrategy struct {
	Name string
}

func (p *AlwaysPlayFirstCardStrategy) Act(otherPlayersCards map[PlayerIndex][]Card, myNumCards int, blueTokens int, redTokens int) Action {
	return Action{Play: &PlayAction{Index: 0}}
}

func (p *AlwaysPlayFirstCardStrategy) ObserveAction(actor PlayerIndex, action Action) {
	log.Printf("%s observed '%s' (by player %d)\n", p.Name, action.DebugString(), actor)
}

func main() {
	fmt.Println("Hello, world!")
	
	deck := shuffle(createDeck())
	players := []PlayerStrategy{
		&AlwaysPlayFirstCardStrategy{"Matt"},
		&AlwaysPlayFirstCardStrategy{"Cristina"}}

	state, err := initializeGame(deck, players)
	
	if err != nil {
		log.Fatal(err)
	}

	for i := PlayerIndex(0); int(i) < len(players); i++ {
		fmt.Printf("PLAYER %d\n", i)
		displayDeck(state.playerStates[i].cards)
	}

	for {
		log.Println("---")
		state.takeTurn()
	}
}

type Card struct {
	Value Value
	Color Color
}

type playerState struct {
	cards []Card
	strategy PlayerStrategy
}

type gameState struct {
	drawPile []Card
	pileHeights map[Color]int
	playerStates []*playerState

	redTokens int  // bad plays
	blueTokens int  // available information

	currentPlayer PlayerIndex
}

func (game *gameState) drawReplacement(player *playerState, card HandIndex) {
	if len(game.drawPile) > 0 {
		// draw a card
		drawn := game.drawPile[0]
		player.cards[card] = drawn
		game.drawPile = game.drawPile[1:]
		log.Printf("Player %d drew a %s %d\n", game.currentPlayer,
			COLOR_INFOS[drawn.Color].fullName, drawn.Value)	
	} else {
		// nothing to draw, remove this card
		player.cards[card] = player.cards[len(player.cards) - 1]
		player.cards = player.cards[1:]
		log.Printf("Nothing left to draw for player %d\n", game.currentPlayer)
	}
}

func (game *gameState) handleGiveInformationAction(action *GiveInformationAction) {
	if game.blueTokens < 1 {
		panic("Invalid action: not enough blue tokens")
	}

	// TODO(mrjones): validate that GiveInformationAction is valid
	// (i.e. only one of Color and Value is set)

	// TODO(mrjones): validate that the information given is correct and complete
	recipient := game.playerStates[action.PlayerIndex]
	for candidateCardPos, candidateCard := range(recipient.cards) {
		givingInformationAboutThisCard := false
		for _, actualCardPos := range(action.Cards) {
			if actualCardPos == HandIndex(candidateCardPos) {
				givingInformationAboutThisCard = true
			}
		}

		if givingInformationAboutThisCard {
			// check that the information matches this card
			if action.Color != nil && candidateCard.Color != *action.Color {
				panic("GiveInformationAction color does not match actual card color")
			}
			if action.Value != nil && candidateCard.Value != *action.Value {
				panic("GiveInformationAction value does not match actual card value")
			}
		} else {
			// check that the information does NOT apply to this card
			if action.Color != nil && candidateCard.Color == *action.Color {
				panic("GiveInformationAction color matches un-referenced card")
			}
			if action.Value != nil && candidateCard.Value == *action.Value {
				panic("GiveInformationAction value matches un-referenced card")
			}
			
		}
	}

	game.blueTokens--
	
}

func (game *gameState) handleDiscardAction(player *playerState, action *DiscardAction) {
	log.Println("TODO: implement discard")
	// TODO(mrjones): check bounds
	card := player.cards[action.Index]
	log.Printf("Player %d discards a %s %d\n",
		game.currentPlayer, COLOR_INFOS[card.Color].fullName, card.Value)

	game.drawReplacement(player, action.Index)

	if game.blueTokens < MAX_BLUE_TOKENS {
		game.blueTokens++
	}
}

func (game *gameState) handlePlayAction(player *playerState, action *PlayAction) {
	// TODO(mrjones): check bounds
	card := player.cards[action.Index]
	log.Printf("Player %d plays a %s %d\n",
		game.currentPlayer, COLOR_INFOS[card.Color].fullName, card.Value)
	if int(card.Value) ==  game.pileHeights[card.Color] + 1 {
		// successful play
		game.pileHeights[card.Color]++
		log.Printf("Good play! %s pile now has height: %d\n",
			COLOR_INFOS[card.Color].fullName, game.pileHeights[card.Color])
		// TODO(mrjones): check if we won the game
	} else {
		// unsuccessful play
		game.redTokens--
		log.Printf("That was a bad play. Red tokens left: %d\n", game.redTokens)
		if game.redTokens == 0 {
			log.Println("We lost the game!")
			panic("we lost the game")
			// TODO(mrjones): we lost the game
		}
	}

	game.drawReplacement(player, action.Index)
}

func (game *gameState) takeTurn() {
	player := game.playerStates[game.currentPlayer]

	otherPlayersCards := make(map[PlayerIndex][]Card)

	for i, player := range(game.playerStates) {
		if PlayerIndex(i) != game.currentPlayer {
			otherPlayersCards[PlayerIndex(i)] = player.cards
		}
	}

	action := player.strategy.Act(
		otherPlayersCards, len(player.cards), game.blueTokens, game.redTokens)
	// TODO(mrjones): validate action

	log.Printf("Player %d taking action '%s'\n", game.currentPlayer, action.DebugString())

	switch {
	case action.GiveInformation != nil:
		game.handleGiveInformationAction(action.GiveInformation)
	case action.Discard != nil:
		game.handleDiscardAction(player, action.Discard)
	case action.Play != nil:
		game.handlePlayAction(player, action.Play)
	default:
		panic("INVALID ACTION")
	}

	for i, player := range(game.playerStates) {
		if PlayerIndex(i) != game.currentPlayer {
			player.strategy.ObserveAction(game.currentPlayer, action)
		}
	}


	game.currentPlayer = PlayerIndex(
		(int(game.currentPlayer) + 1) % len(game.playerStates))
}

func initializeGame(deck []Card, players []PlayerStrategy) (*gameState, error) {
	numPlayers := len(players)

	cardsPerPlayer, ok := INITIAL_CARDS[numPlayers]
	if !ok {
		return nil, fmt.Errorf("Invalid number of players: %d", numPlayers)
	}

	state := &gameState{
		playerStates: make([]*playerState, numPlayers),
		pileHeights: make(map[Color]int),
		drawPile: []Card{},
		currentPlayer: PlayerIndex(rand.Intn(numPlayers)),
		redTokens: 3,
		blueTokens: MAX_BLUE_TOKENS,
	}

	for i, _ := range(ALL_COLORS) {
		state.pileHeights[Color(i)] = 0;
	}

	for p := PlayerIndex(0); int(p) < numPlayers; p++ {
		state.playerStates[p] = &playerState{
			cards: []Card{},
			strategy: players[p],
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
