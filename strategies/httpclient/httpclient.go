package httpclient

import (
	"github.com/mrjones/yanhuo/core"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpClientStrategy struct {
	remoteEndpoint *url.URL
	httpClient *http.Client
}

func NewHttpClientStrategy(remoteEndpoint *url.URL) *HttpClientStrategy {
	return &HttpClientStrategy {
		remoteEndpoint: remoteEndpoint,
		httpClient: &http.Client{},
	}
}

type GameState struct {
	MyPlayerIndex int
	OtherPlayersCards map[string][]yanhuo.Card
	MyCardCount int
	BlueTokens int
	RedTokens int
}

type Observation struct {
	Actor yanhuo.PlayerIndex
	Action yanhuo.Action
}

type Transmission struct {
	MessageType string
	GameState *GameState `json:",omitempty"`
	Observation *Observation `json:",omitempty"`
}

func translateCardMap(in map[yanhuo.PlayerIndex][]yanhuo.Card) map[string][]yanhuo.Card {
	out := map[string][]yanhuo.Card{}

	for idx, cards := range(in) {
		out[fmt.Sprintf("%d", idx)] = cards
	}

	return out
}

func (p *HttpClientStrategy) StartGame(
	myPlayerIndex yanhuo.PlayerIndex,
	otherPlayersCards map[yanhuo.PlayerIndex][]yanhuo.Card,
	myNumCards int,
	blueTokens int,
	redTokens int) {
	transmission := Transmission{
		MessageType: "StartGame",
		GameState: &GameState{
			MyPlayerIndex: int(myPlayerIndex),
			OtherPlayersCards: translateCardMap(otherPlayersCards),
			MyCardCount: myNumCards,
			BlueTokens: blueTokens,
			RedTokens: redTokens,
		},
	}
	
	payload, err := json.Marshal(transmission)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transmitting %s\n", string(payload))
	_, err = p.httpClient.Post(p.remoteEndpoint.String(), "application/json", bytes.NewReader(payload))
	
	if err != nil {
		panic(err)
	}
}

func (p *HttpClientStrategy) Act(
	myPlayerIndex yanhuo.PlayerIndex,
	otherPlayersCards map[yanhuo.PlayerIndex][]yanhuo.Card,
	myNumCards int,
	blueTokens int,
	redTokens int) yanhuo.Action {
	transmission := Transmission{
		MessageType: "ActionRequest",
		GameState: &GameState{
			MyPlayerIndex: int(myPlayerIndex),
			OtherPlayersCards: translateCardMap(otherPlayersCards),
			MyCardCount: myNumCards,
			BlueTokens: blueTokens,
			RedTokens: redTokens,
		},
	}
	
	payload, err := json.Marshal(transmission)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transmitting %s\n", string(payload))
	resp, err := p.httpClient.Post(p.remoteEndpoint.String(), "application/json", bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}

	var action yanhuo.Action

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Unmarshaling: %s\n", string(respBytes))

	err = json.Unmarshal(respBytes, &action)
	if err != nil {
		panic(err)
	}

	return action
}

func (p *HttpClientStrategy) ObserveAction(actor yanhuo.PlayerIndex, action yanhuo.Action) {
	transmission := Transmission{
		MessageType: "Observation",
		Observation: &Observation{Actor: actor, Action: action},
	}

	payload, err := json.Marshal(transmission)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transmitting %s\n", string(payload))
	p.httpClient.Post(p.remoteEndpoint.String(), "application/json", bytes.NewReader(payload))
}

