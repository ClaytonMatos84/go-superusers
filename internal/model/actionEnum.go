package model

import (
	"encoding/json"
	"fmt"
)

type Action int

const (
	Login Action = iota
	Logout
)

func (a Action) String() string {
	switch a {
	case Login:
		return "login"
	case Logout:
		return "logout"
	default:
		return "unknown"
	}
}

func (a Action) MarshalJSON() ([]byte, error) {
	switch a {
	case Login:
		return []byte(`"login"`), nil
	case Logout:
		return []byte(`"logout"`), nil
	default:
		return nil, nil
	}
}
func (a *Action) UnmarshalJSON(b []byte) error {
	var actionStr string
	if err := json.Unmarshal(b, &actionStr); err != nil {
		return err
	}

	switch actionStr {
	case "login":
		*a = Login
	case "logout":
		*a = Logout
	default:
		return fmt.Errorf("invalid action: %s", actionStr)
	}
	return nil
}
