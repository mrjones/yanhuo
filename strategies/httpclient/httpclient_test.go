package httpclient

import (
	"github.com/mrjones/yanhuo/core"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type TestRoundTripper struct {
	LastRequest *http.Request
	Response *http.Response
}

func (tr *TestRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	tr.LastRequest = request
	return tr.Response, nil
}

func parseTransmission(t *testing.T, req *http.Request) (*Transmission) {
	reqBodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Req body: %s\n", string(reqBodyBytes))
	
	var transmission *Transmission
	err = json.Unmarshal(reqBodyBytes, &transmission)
	if err != nil {
		t.Fatal(err)
	}

	return transmission
}

func makeStrategy(t *testing.T) (*HttpClientStrategy, *TestRoundTripper) {
	rt := &TestRoundTripper{}
	u, err := url.Parse("http://www.example.com/game")
	if err != nil {
		t.Fatal(err)
	}
	s := NewHttpClientStrategy(u)

	s.httpClient = &http.Client{Transport: rt}

	return s, rt
}

func makeOkResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(strings.NewReader(body)),
	}
}

func TestObserve(t *testing.T) {
	s, rt := makeStrategy(t)
	rt.Response = makeOkResponse("")

	a := yanhuo.Action{
		Discard: &yanhuo.DiscardAction{Index: yanhuo.HandIndex(0)},
	}
	s.ObserveAction(yanhuo.PlayerIndex(2), a)

	transmission := parseTransmission(t, rt.LastRequest)

	if transmission.Observation.Actor != yanhuo.PlayerIndex(2) {
		t.Errorf("Expected player (2) does not match actual: %d",
			transmission.Observation.Actor)
	}

	if a.DebugString() != transmission.Observation.Action.DebugString() {
		t.Errorf("Original object:\n'%s'\nDoes not match object parsed on server:\n'%s'", a.DebugString(), transmission.Observation.Action.DebugString())
	}
}

func makeCard(value int, color yanhuo.Color) yanhuo.Card {
	return yanhuo.Card{
		Value: yanhuo.Value(value),
		Color: color,
	}
}

func TestAct(t *testing.T) {
	s, rt := makeStrategy(t)
	rt.Response = makeOkResponse("{\"Discard\":{\"Index\":2}}")

	decision := s.Act(
		yanhuo.PlayerIndex(2),
		map[yanhuo.PlayerIndex][]yanhuo.Card{
			yanhuo.PlayerIndex(0): []yanhuo.Card{
				makeCard(1, yanhuo.RED),
				makeCard(2, yanhuo.RED),
			},
			yanhuo.PlayerIndex(1): []yanhuo.Card{
				makeCard(3, yanhuo.BLUE),
				makeCard(4, yanhuo.YELLOW),
			},
		},
		2, 3, 4)

	tran := parseTransmission(t, rt.LastRequest)

	if tran.MessageType != "ActionRequest" {
		t.Errorf("Wrong transmission.MessageType: %s", tran.MessageType)
	}

	if tran.GameState == nil {
		t.Errorf("Server should have interpreted action request")
	}

	if tran.GameState.MyCardCount != 2 {
		t.Errorf("Wrong card count in transmission: %d", tran.GameState.MyCardCount)
	}

	if tran.GameState.BlueTokens != 3 {
		t.Errorf("Wrong blue count in transmission: %d", tran.GameState.BlueTokens)
	}

	if tran.GameState.RedTokens != 4 {
		t.Errorf("Wrong red count in transmission: %d", tran.GameState.RedTokens)
	}

	if decision.Discard == nil {
		t.Errorf("Should have discarded: %s", decision.DebugString())
	}

	if decision.Discard.Index != 2 {
		t.Errorf("Should have discarded card 2: %s", decision.DebugString())
	}
}
