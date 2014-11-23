// To play a game:
// 1. Call hanabi.InitializeGame() passing in an array of PlayerStrategies
//    (See below for how to implement a PlayerStrategy)
// 2. Repeatedly call TakeTurn on the returned object
//    TODO(mrjones): define when to *stop* calling TakeTurn

package hanabi

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

type Card struct {
	Value Value
	Color Color
}

//
// Interface for implementing new strategies
//

type PlayerStrategy interface {
	Act(
		otherPlayersCards map[PlayerIndex][]Card,
		myNumCards int,
		blueTokens int,
		redTokens int) Action

	ObserveAction(actor PlayerIndex, action Action)
}

type Action struct {
	// Exactly one one must be non-null
	GiveInformation *GiveInformationAction
	Discard *DiscardAction
	Play *PlayAction
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

//
// IMPLEMENTATION DETAILS BELOW HERE
//

const (
	kMaxBlueTokens = 8
)

type colorInfo struct {
	fullName string
	shortName string
}

var kColorInfos = map[Color]colorInfo {
	WHITE: colorInfo{fullName: "WHITE", shortName: "W"},
	RED: colorInfo{fullName: "RED", shortName: "R"},
	BLUE: colorInfo{fullName: "BLUE", shortName: "B"},
	YELLOW: colorInfo{fullName: "YELLOW", shortName: "Y"},
	GREEN: colorInfo{fullName: "GREEN", shortName: "G"},
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
			kColorInfos[drawn.Color].fullName, drawn.Value)	
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
		game.currentPlayer, kColorInfos[card.Color].fullName, card.Value)

	game.drawReplacement(player, action.Index)

	if game.blueTokens < kMaxBlueTokens {
		game.blueTokens++
	}
}

func (game *gameState) handlePlayAction(player *playerState, action *PlayAction) {
	// TODO(mrjones): check bounds
	card := player.cards[action.Index]
	log.Printf("Player %d plays a %s %d\n",
		game.currentPlayer, kColorInfos[card.Color].fullName, card.Value)
	if int(card.Value) ==  game.pileHeights[card.Color] + 1 {
		// successful play
		game.pileHeights[card.Color]++
		log.Printf("Good play! %s pile now has height: %d\n",
			kColorInfos[card.Color].fullName, game.pileHeights[card.Color])
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

func (game *gameState) TakeTurn() {
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

func InitializeGame(players []PlayerStrategy) (*gameState, error) {
	deck := shuffle(createDeck())

	numPlayers := len(players)

	// map from number of players to number of initial cards per player
	var kInitialCards = map[int]int {
		2: 5,
		3: 5,
		4: 4,
		5: 4,
	}

	cardsPerPlayer, ok := kInitialCards[numPlayers]
	if !ok {
		return nil, fmt.Errorf("Invalid number of players: %d", numPlayers)
	}

	state := &gameState{
		playerStates: make([]*playerState, numPlayers),
		pileHeights: make(map[Color]int),
		drawPile: []Card{},
		currentPlayer: PlayerIndex(rand.Intn(numPlayers)),
		redTokens: 3,
		blueTokens: kMaxBlueTokens,
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

func DisplayDeck(deck []Card) {
	for _, card := range(deck) {
		fmt.Printf("color: %s, value: %d\n", kColorInfos[card.Color].fullName, card.Value)
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
	var kNumCardsInDeckByValue = map[Value]int{
		1: 3,
		2: 2,
		3: 2,
		4: 2,
		5: 1,
	}

	cards := []Card{}
	for color, _ := range(kColorInfos) {
		for value, numCardsInDeckWithValue := range(kNumCardsInDeckByValue) {
			for i := 0; i < numCardsInDeckWithValue; i++ {
				cards = append(cards, Card{Value: value, Color: color})
			}
		}
	}
	
	return cards
}
