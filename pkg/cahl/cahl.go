package cahl

import (
	"fmt"
	"log/slog"
)

type Team struct {
	Name    string    `json:"name"`
	Manager string    `json:"manager"`
	Players []*Player `json:"players"`
	Clubs   []*Club   `json:"teams"`
}

type Player struct {
	Name     string      `json:"name"`
	ID       uint64      `json:"id"`
	Position Position    `json:"position"`
	Stats    PlayerStats `json:"stats"`
}

type PlayerStats struct {
	Goals   int `json:"goals"`
	Assists int `json:"assists"`
}

type Club struct {
	Name   string    `json:"name"`
	Abbrev string    `json:"abbrev"`
	Stats  ClubStats `json:"stats"`
}

type ClubStats struct {
	Wins       int `json:"wins"`
	LossesInOT int `json:"losses_in_ot"`
}

// Returns nil if the team is valid or an error if it's not
func (t Team) Valid() error {
	if len(t.Players) != 9 {
		return fmt.Errorf("team '%s' has the wrong number of players (%d)", t.Name, len(t.Players))
	}

	if len(t.Clubs) != 3 {
		return fmt.Errorf("team '%s' has the wrong number of clubs (%d)", t.Name, len(t.Clubs))
	}

	return nil
}

func (t Team) Score() (score int) {
	for _, p := range t.Players {
		score += p.Score()
	}

	for _, c := range t.Clubs {
		score += c.Score()
	}

	slog.Debug("calculated team score", "team", t.Name, "score", score)

	return
}

func (p Player) Score() (score int) {
	if p.Position == Defence {
		score += p.Stats.Goals * 3
	} else {
		score += p.Stats.Goals * 2
	}

	score += p.Stats.Assists

	slog.Debug("player score", "name", p.Name, "score", score, "assists", p.Stats.Assists, "goals", p.Stats.Goals, "position", p.Position)

	return
}

func (c Club) Score() (score int) {
	slog.Debug("club score", "name", c.Abbrev, "wins", c.Stats.Wins, "lossesInOT", c.Stats.LossesInOT, "inc", c.Stats.Wins*2+c.Stats.LossesInOT)
	score += c.Stats.Wins*2 + c.Stats.LossesInOT

	return
}
