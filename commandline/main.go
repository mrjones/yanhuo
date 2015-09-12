package main

import (
	"github.com/mrjones/yanhuo/core"
	"github.com/mrjones/yanhuo/strategies/alwaysplay"

	"log"
)

func main() {
	state, err := yanhuo.InitializeGame(
		[]yanhuo.PlayerStrategy{
			&alwaysplay.AlwaysPlayFirstCardStrategy{"Matt"},
			&alwaysplay.AlwaysPlayFirstCardStrategy{"Cristina"},
		},
		[]yanhuo.Observer{
			&yanhuo.LoggingObserver{},
		})

	if err != nil {
		log.Fatal(err)
	}

	state.Play()
}
