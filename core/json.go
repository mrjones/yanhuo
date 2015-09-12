package yanhuo

import (
	"encoding/json"
	"fmt"
)

func (ci *ColorInformation) MarshalJSON() ([]byte, error) {
	return json.Marshal(kColorInfos[ci.Color].fullName)
}

func (ci *ColorInformation) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("color should be a string, got %s", data)
	}

	for color, colorInfo := range(kColorInfos) {
		if s == colorInfo.fullName {
			fmt.Printf("Matched: %d %s %s\n", color, colorInfo.fullName, colorInfo.shortName)
			ci.Color = color
			return nil
		}
	}

	return fmt.Errorf("invalid color %q", s)
}

func (vi *ValueInformation) MarshalJSON() ([]byte, error) {
	return json.Marshal(vi.Value)
}

func (vi *ValueInformation) UnmarshalJSON(data []byte) error {
	var v Value
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("value should be a int, got %s", data)
	}

	vi.Value = v
	return nil
}
