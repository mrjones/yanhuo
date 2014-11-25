package main

import (
	"github.com/mrjones/yanhuo/yanhuo"
	"github.com/mrjones/yanhuo/strategies"

	"log"
)

func main() {
	state, err := yanhuo.InitializeGame(
		[]yanhuo.PlayerStrategy{
			&strategies.AlwaysPlayFirstCardStrategy{"Matt"},
			&strategies.AlwaysPlayFirstCardStrategy{"Cristina"},
		},
		[]yanhuo.Observer{
			&yanhuo.LoggingObserver{},
		})

	if err != nil {
		log.Fatal(err)
	}

	state.Play()
}
