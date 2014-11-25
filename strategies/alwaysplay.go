package strategies

import (
	"github.com/mrjones/yanhuo/yanhuo"

	"log"
)

// TODO(mrjones): move to a separate package
type AlwaysPlayFirstCardStrategy struct {
	Name string
}

func (p *AlwaysPlayFirstCardStrategy) Act(
	otherPlayersCards map[yanhuo.PlayerIndex][]yanhuo.Card,
	myNumCards int,
	blueTokens int,
	redTokens int) yanhuo.Action {
	return yanhuo.Action{Play: &yanhuo.PlayAction{Index: 0}}
}

func (p *AlwaysPlayFirstCardStrategy) ObserveAction(actor yanhuo.PlayerIndex, action yanhuo.Action) {
	log.Printf("%s observed '%s' (by player %d)\n", p.Name, action.DebugString(), actor)
}
