package cahl

import (
	"encoding/json"
	"errors"
)

type Position int

const (
	Unknown Position = iota
	Forward
	Defence
)

func (p Position) String() string {
	switch p {
	case Forward:
		return "forward"
	case Defence:
		return "defence"
	default:
		return "unknown"
	}
}

func ParsePositionFromAPI(s string) (Position, error) {
	switch s {
	case "C":
		fallthrough
	case "R":
		fallthrough
	case "L":
		return Forward, nil
	case "D":
		return Defence, nil
	default:
		return Unknown, errors.New("unknown position")
	}
}

func ParseMarshaledPosition(s string) (Position, error) {
	switch s {
	case "forward":
		return Forward, nil
	case "defence":
		return Defence, nil
	}

	return Unknown, errors.New("unknown position")
}

func (p Position) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Position) UnmarshalJSON(data []byte) (err error) {
	var pos string

	if err := json.Unmarshal(data, &pos); err != nil {
		return err
	}

	if *p, err = ParseMarshaledPosition(pos); err != nil {
		return err
	}

	return nil
}
