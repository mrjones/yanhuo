package main

import (
	"github.com/mrjones/yanhuo/core"
	"github.com/mrjones/yanhuo/strategies/alwaysplay"
	"github.com/mrjones/yanhuo/strategies/httpclient"

	"log"
	"net/url"
)

func main() {
	remoteUrl, err := url.Parse("http://192.168.1.4:8000/")
	if err != nil {
		log.Fatal(err)
	}

	state, err := yanhuo.InitializeGame(
		[]yanhuo.PlayerStrategy{
			&alwaysplay.AlwaysPlayFirstCardStrategy{"Matt"},
			&alwaysplay.AlwaysPlayFirstCardStrategy{"Cristina"},
			httpclient.NewHttpClientStrategy(remoteUrl),
		},
		[]yanhuo.Observer{
			&yanhuo.LoggingObserver{},
		})

	if err != nil {
		log.Fatal(err)
	}

	state.Play()
}
