// To play a game:
// 1. Call yanhuo.InitializeGame() passing in an array of PlayerStrategies
//    (See below for how to implement a PlayerStrategy)
// 2. Call Play on the returned object

package yanhuo

import (
	"bytes"
	"fmt"
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
	StartGame(
		myPlayerIndex PlayerIndex,
		otherPlayersCards map[PlayerIndex][]Card,
		myNumCards int,
		blueTokens int,
		redTokens int)

	Act(
		myPlayerIndex PlayerIndex,
		otherPlayersCards map[PlayerIndex][]Card,
		myNumCards int,
		blueTokens int,
		redTokens int) Action

	ObserveAction(actor PlayerIndex, action Action)
}

type Action struct {
	// Exactly one one must be non-null
	GiveInformation *GiveInformationAction `json:",omitempty"`
	Discard         *DiscardAction `json:",omitempty"`
	Play            *PlayAction `json:",omitempty"`
}

type ColorInformation struct {
	Color Color
}

type ValueInformation struct {
	Value Value
}

type GiveInformationAction struct {
	// The player information is being given about
	PlayerIndex PlayerIndex

	// The cards matching the information being given
	Cards []HandIndex

	// Exactly one of color or value must be non-null
	Color *ColorInformation `json:",omitempty"`
	Value *ValueInformation `json:",omitempty"`
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
// Interface for observing games
//

type Observer interface {
	GameStart(cards [][]Card)

	ObserveAction(p PlayerIndex, a Action)
	ObserveDiscard(p PlayerIndex, c Card, i HandIndex)
	ObserveDraw(p PlayerIndex, c Card, i HandIndex)
	ObservePlay(p PlayerIndex, c Card, successful bool)

	TurnComplete(piles map[Color]int, blueTokens int, redTokens int)
	GameComplete(won bool, piles map[Color]int)
}

//
// IMPLEMENTATION DETAILS BELOW HERE
//

const (
	kMaxBlueTokens = 8

	kKeepGoing = true
	kStop      = false
)

type colorInfo struct {
	fullName  string
	shortName string
}

var kColorInfos = map[Color]colorInfo{
	WHITE:  colorInfo{fullName: "WHITE", shortName: "W"},
	RED:    colorInfo{fullName: "RED", shortName: "R"},
	BLUE:   colorInfo{fullName: "BLUE", shortName: "B"},
	YELLOW: colorInfo{fullName: "YELLOW", shortName: "Y"},
	GREEN:  colorInfo{fullName: "GREEN", shortName: "G"},
}

func (a Action) QuickString() string {
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

const kValid = ""

func (a Action) IsValid() bool {
	return a.InvalidReason() == kValid
}

func (a Action) InvalidReason() string {
	subActions := 0
	if a.GiveInformation != nil {
		subActions++
		subReason := a.GiveInformation.invalidReason()
		if subReason != kValid {
			return subReason
		}
	}
	if a.Discard != nil {
		subActions++
	}
	if a.Play != nil {
		subActions++
	}

	if subActions == 1 {
		return kValid
	} else if subActions == 0 {
		return "No sub action (GiveInformation, Play, Discard) was set."
	} else {
		return "More than one sub action (GiveInformation, Play, Discard) was set."
	}
}

func (a *GiveInformationAction) invalidReason() string {
	informationTypes := 0
	if a.Color != nil {
		informationTypes++
	}
	if a.Value != nil {
		informationTypes++
	}

	if informationTypes == 1 {
		return kValid
	} else if informationTypes == 0 {
		return "No information type (Color, Value) was set."
	} else {
		return "More than one information type (Color, Value) was set."
	}

	return kValid
}

type playerState struct {
	cards    []Card
	strategy PlayerStrategy
}

type gameState struct {
	drawPile     []Card
	pileHeights  map[Color]int
	playerStates []*playerState
	observers    []Observer

	redTokens  int // bad plays
	blueTokens int // available information

	currentPlayer PlayerIndex

	finished bool
	won      bool
}

func (game *gameState) drawReplacement(player *playerState, card HandIndex) {
	if len(game.drawPile) > 0 {
		// draw a card
		drawn := game.drawPile[0]
		player.cards[card] = drawn
		game.drawPile = game.drawPile[1:]
		for _, o := range game.observers {
			o.ObserveDraw(game.currentPlayer, drawn, card)
		}
	} else {
		// nothing to draw, remove this card
		player.cards[card] = player.cards[len(player.cards)-1]
		player.cards = player.cards[1:]
	}
}

func (game *gameState) handleGiveInformationAction(action *GiveInformationAction) bool {
	if game.blueTokens < 1 {
		panic("Invalid action: not enough blue tokens")
	}

	if action.invalidReason() != kValid {
		panic("Invalid action:" + action.invalidReason())
	}

	recipient := game.playerStates[action.PlayerIndex]
	for candidateCardPos, candidateCard := range recipient.cards {
		givingInformationAboutThisCard := false
		for _, actualCardPos := range action.Cards {
			if actualCardPos == HandIndex(candidateCardPos) {
				givingInformationAboutThisCard = true
			}
		}

		if givingInformationAboutThisCard {
			// check that the information matches this card
			if action.Color != nil && candidateCard.Color != action.Color.Color {
				panic("Invalid action: GiveInformationAction color does not match actual card color")
			}
			if action.Value != nil && candidateCard.Value != action.Value.Value {
				panic("Invalid action: GiveInformationAction value does not match actual card value")
			}
		} else {
			// check that the information does NOT apply to this card
			if action.Color != nil && candidateCard.Color == action.Color.Color {
				panic("Invalid action: GiveInformationAction color matches un-referenced card")
			}
			if action.Value != nil && candidateCard.Value == action.Value.Value {
				panic("Invalid action: GiveInformationAction value matches un-referenced card")
			}

		}
	}

	game.blueTokens--
	return kKeepGoing
}

func (game *gameState) handleDiscardAction(player *playerState, action *DiscardAction) bool {
	if int(action.Index) >= len(player.cards) {
		panic(fmt.Sprintf("Invalid action: Index (%d) was out of bounds (len: %d)",
			action.Index, len(player.cards)))
	}

	card := player.cards[action.Index]
	for _, o := range game.observers {
		o.ObserveDiscard(game.currentPlayer, card, action.Index)
	}

	game.drawReplacement(player, action.Index)

	if game.blueTokens < kMaxBlueTokens {
		game.blueTokens++
	}

	return kKeepGoing
}

func (game *gameState) handlePlayAction(player *playerState, action *PlayAction) bool {
	if int(action.Index) >= len(player.cards) {
		panic(fmt.Sprintf("Invalid action: Index (%d) was out of bounds (len: %d)",
			action.Index, len(player.cards)))
	}

	card := player.cards[action.Index]
	success := int(card.Value) == game.pileHeights[card.Color]+1
	for _, o := range game.observers {
		o.ObservePlay(game.currentPlayer, card, success)
	}
	if success {
		// successful play
		game.pileHeights[card.Color]++
		won := true
		for _, color := range ALL_COLORS {
			if game.pileHeights[color] != 5 {
				won = false
			}
		}
		if won {
			game.finished = true
			game.won = true
			return kStop
		}
	} else {
		// unsuccessful play
		game.redTokens--
		if game.redTokens == 0 {
			game.finished = true
			game.won = false
			return kStop
		}
	}

	game.drawReplacement(player, action.Index)
	return kKeepGoing
}

func (game *gameState) Play() bool {

	for i, player := range game.playerStates {
		// TODO(mrjones): factor out duplicated code
		otherPlayersCards := make(map[PlayerIndex][]Card)
		
		for j, player := range game.playerStates {
			if i != j {
				otherPlayersCards[PlayerIndex(j)] = player.cards
			}
		}

		player.strategy.StartGame(
			PlayerIndex(i), otherPlayersCards, len(player.cards), game.blueTokens, game.redTokens)
	}

	for game.takeTurn() {
	}

	for _, o := range game.observers {
		o.GameComplete(game.won, game.pileHeights)
	}

	return game.won
}

func (game *gameState) takeTurn() bool {
	player := game.playerStates[game.currentPlayer]

	// TODO(mrjones): factor out duplicated code
	otherPlayersCards := make(map[PlayerIndex][]Card)

	for i, player := range game.playerStates {
		if PlayerIndex(i) != game.currentPlayer {
			otherPlayersCards[PlayerIndex(i)] = player.cards
		}
	}

	action := player.strategy.Act(
		game.currentPlayer, otherPlayersCards, len(player.cards), game.blueTokens, game.redTokens)

	if !action.IsValid() {
		panic("Invalid action: " + action.InvalidReason())
	}

	for _, o := range game.observers {
		o.ObserveAction(game.currentPlayer, action)
	}

	keepGoing := true
	switch {
	case action.GiveInformation != nil:
		keepGoing = game.handleGiveInformationAction(action.GiveInformation)
	case action.Discard != nil:
		keepGoing = game.handleDiscardAction(player, action.Discard)
	case action.Play != nil:
		keepGoing = game.handlePlayAction(player, action.Play)
	default:
		panic("INVALID ACTION")
	}

	for i, player := range game.playerStates {
		if PlayerIndex(i) != game.currentPlayer {
			player.strategy.ObserveAction(game.currentPlayer, action)
		}
	}

	game.currentPlayer = PlayerIndex(
		(int(game.currentPlayer) + 1) % len(game.playerStates))

	for _, o := range game.observers {
		o.TurnComplete(game.pileHeights, game.blueTokens, game.redTokens)
	}

	return keepGoing
}

func InitializeGame(players []PlayerStrategy, observers []Observer) (*gameState, error) {
	deck := shuffle(createDeck())

	numPlayers := len(players)

	// map from number of players to number of initial cards per player
	var kInitialCards = map[int]int{
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
		playerStates:  make([]*playerState, numPlayers),
		pileHeights:   make(map[Color]int),
		drawPile:      []Card{},
		currentPlayer: PlayerIndex(rand.Intn(numPlayers)),
		redTokens:     3,
		blueTokens:    kMaxBlueTokens,
		observers:     observers,
		finished:      false,
		won:           false,
	}

	for i, _ := range ALL_COLORS {
		state.pileHeights[Color(i)] = 0
	}

	for p := PlayerIndex(0); int(p) < numPlayers; p++ {
		state.playerStates[p] = &playerState{
			cards:    []Card{},
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

	for _, o := range state.observers {
		cards := make([][]Card, len(state.playerStates))
		for i, player := range state.playerStates {
			cards[i] = player.cards
		}

		o.GameStart(cards)
	}

	state.drawPile = deck[drawCount:]
	return state, nil
}

func DisplayDeck(deck []Card) {
	for _, card := range deck {
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
	for color, _ := range kColorInfos {
		for value, numCardsInDeckWithValue := range kNumCardsInDeckByValue {
			for i := 0; i < numCardsInDeckWithValue; i++ {
				cards = append(cards, Card{Value: value, Color: color})
			}
		}
	}

	return cards
}

func (a *Action) DebugString() string {
	var buffer bytes.Buffer

	if a.GiveInformation != nil {
		buffer.WriteString("GiveInformation < ")
		buffer.WriteString(a.GiveInformation.DebugString())
		buffer.WriteString("> ")
	}

	if a.Discard != nil {
		buffer.WriteString("Discard < ")
		buffer.WriteString(a.Discard.DebugString())
		buffer.WriteString("> ")
	}

	if a.Play != nil {
		buffer.WriteString("Play < ")
		buffer.WriteString(a.Play.DebugString())
		buffer.WriteString("> ")
	}
	
	return buffer.String()
}

func (a *GiveInformationAction) DebugString() string {
	var buffer bytes.Buffer 

	buffer.WriteString(fmt.Sprintf("PlayerIndex:%d ", a.PlayerIndex))
	buffer.WriteString(fmt.Sprintf("Cards:%v ", a.Cards))

	if a.Color != nil {
		buffer.WriteString(fmt.Sprintf(
			"Color:%s ", kColorInfos[a.Color.Color].fullName))
	}

	if a.Value != nil {
		buffer.WriteString(fmt.Sprintf("Value:%d ", a.Value.Value))
	}

	return buffer.String()
}

func (a *DiscardAction) DebugString() string {
	return fmt.Sprintf("Index:%d ", a.Index)
}

func (a *PlayAction) DebugString() string {
	return fmt.Sprintf("Index:%d ", a.Index)
}
