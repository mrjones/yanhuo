package strategies

import (
	"github.com/mrjones/hanabi/hanabi"

	"log"
)

// TODO(mrjones): move to a separate package
type AlwaysPlayFirstCardStrategy struct {
	Name string
}

func (p *AlwaysPlayFirstCardStrategy) Act(
	otherPlayersCards map[hanabi.PlayerIndex][]hanabi.Card,
	myNumCards int,
	blueTokens int,
	redTokens int) hanabi.Action {
	return hanabi.Action{Play: &hanabi.PlayAction{Index: 0}}
}

func (p *AlwaysPlayFirstCardStrategy) ObserveAction(actor hanabi.PlayerIndex, action hanabi.Action) {
	log.Printf("%s observed '%s' (by player %d)\n", p.Name, action.DebugString(), actor)
}
