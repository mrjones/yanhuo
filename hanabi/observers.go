package hanabi

import (
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
	log.Printf("Turn complete.\n---\n")
}
