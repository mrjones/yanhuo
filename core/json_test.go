package yanhuo

import (
	"encoding/json"
	"testing"
)

func roundTripAction(t *testing.T, a1 *Action, s string) {
	data, err := json.Marshal(a1)
	if err != nil {
		t.Error(err)
		return
	}

	if string(data) != s {
		t.Errorf("Expected JSON:\n%s\nActual JSON:\n%s", s, string(data))
	}

	var a2 *Action
	err = json.Unmarshal(data, &a2)
	if err != nil {
		t.Error(err)
		return
	}

	a1s := a1.DebugString()
	a2s := a2.DebugString()

	if a1s != a2s {
		t.Errorf("Original object:\n'%s'\nDoes not match round-tripped object:\n'%s'", a1.DebugString(), a2.DebugString())
	}
}

func TestRoundTrips_Action(t *testing.T) {
	roundTripAction(t, &Action{
		GiveInformation: &GiveInformationAction{
			PlayerIndex: 2,
			Cards: []HandIndex{1, 3},
			Color: &ColorInformation{Color: RED},
		},
	}, "{\"GiveInformation\":{\"PlayerIndex\":2,\"Cards\":[1,3],\"Color\":\"RED\"}}")

	roundTripAction(t, &Action{
		GiveInformation: &GiveInformationAction{
			PlayerIndex: 2,
			Cards: []HandIndex{1, 3},
			Value: &ValueInformation{Value: 4},
		},
	}, "{\"GiveInformation\":{\"PlayerIndex\":2,\"Cards\":[1,3],\"Value\":4}}")

	roundTripAction(t, &Action{
		Play: &PlayAction{Index: 3},
	}, "{\"Play\":{\"Index\":3}}")

	roundTripAction(t, &Action{
		Discard: &DiscardAction{Index: 3},
	}, "{\"Discard\":{\"Index\":3}}")
}


