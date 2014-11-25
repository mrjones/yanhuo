package hanabi

import (
	"fmt"
	"log"
)

type LoggingObserver struct {
}

func (o *LoggingObserver) ObserveAction(p PlayerIndex, a Action) {
	log.Printf("Player %d taking action '%s'\n", p, a.DebugString())
}

func (o *LoggingObserver) ObserveDiscard(p PlayerIndex, c Card, i HandIndex) {
	log.Printf("Player %d discards a %s %d\n", p, kColorInfos[c.Color].fullName, c.Value)
}

func (o *LoggingObserver) ObserveDraw(p PlayerIndex, c Card, i HandIndex) {
	log.Printf("Player %d drew a %s %d\n", p, kColorInfos[c.Color].fullName, c.Value)
}

func (o *LoggingObserver) ObservePlay(p PlayerIndex, c Card, successful bool) {
	log.Printf("Player %d played a %s %d. Successful: %t\n", p,
		kColorInfos[c.Color].fullName, c.Value, successful)
}

func (o *LoggingObserver) TurnComplete(piles map[Color]int, blueTokens int, redTokens int) {
	log.Printf("Turn complete.\n")
	log.Printf("---\n")
}

func (o *LoggingObserver) GameComplete(won bool, piles map[Color]int) {
	if won {
		log.Printf("We won!")
	} else {
		log.Printf("We lost.")
	}
}

func (o *LoggingObserver) GameStart(cards [][]Card) {
	for i, playerCards := range cards {
		log.Printf("Player %d: [%s]\n", i, summarizeCards(playerCards))
	}
	log.Printf("===\n")
}

func summarizeCards(cards []Card) string {
	a := ""
	sep := ""

	for _, c := range cards {
		a = fmt.Sprintf("%s%s%s%d", a, sep, kColorInfos[c.Color].shortName, c.Value)
		sep = ", "
	}

	return a
}
