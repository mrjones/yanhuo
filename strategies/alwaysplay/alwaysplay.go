package alwaysplay

import (
	"github.com/mrjones/yanhuo/core"

	"log"
)

type AlwaysPlayFirstCardStrategy struct {
	Name string
}

func (p *AlwaysPlayFirstCardStrategy) StartGame(
	myPlayerIndex yanhuo.PlayerIndex,
	otherPlayersCards map[yanhuo.PlayerIndex][]yanhuo.Card,
	myNumCards int,
	blueTokens int,
	redTokens int) { }

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
