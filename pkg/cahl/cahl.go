package cahl

import (
	"fmt"
	"log/slog"
)

type Team struct {
	Name        string    `json:"name"`
	Manager     string    `json:"manager"`
	Players     []*Player `json:"players"`
	Clubs       []*Club   `json:"teams"`
	QuebecRound Player    `json:"quebec_round"`
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
	if len(t.Players) != TeamPlayerCount {
		return fmt.Errorf("team '%s' has the wrong number of players (%d)", t.Name, len(t.Players))
	}

	if len(t.Clubs) != TeamClubCount {
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

const (
	TeamPlayerCount = 9
	TeamClubCount   = 3

	GoalDefencePointsFactor = 3
	GoalForwardPointsFactor = 2
	AssistPointsFactor      = 1
	WinPointsFactor         = 2
	LossesInOTPointsFactor  = 1
	QuebecRoundGoalFactor   = 1
	QuebecRoundMultiplier   = 3
)

func (p Player) ScoreForQuebecRound() (int, error) {
	var quebecPlayerThresholds = map[string]int{
		"Claude Giroux":      14,
		"Philippe Danault":   14,
		"Zachary Bolduc":     14,
		"Alexis Lafrenière":  27,
		"Pierre-Luc Dubois":  21,
		"Jonathan Huberdeau": 23,
	}

	threshold, ok := quebecPlayerThresholds[p.Name]
	if !ok {
		return 0, fmt.Errorf("invalid name '%s'", p.Name)
	}

	score := p.Stats.Goals * QuebecRoundGoalFactor

	if p.Stats.Goals >= threshold {
		score *= QuebecRoundMultiplier
	}

	return score, nil
}

func (p Player) ScoreForGoals() (score int) {
	if p.Position == Defence {
		score += p.Stats.Goals * GoalDefencePointsFactor
	} else {
		score += p.Stats.Goals * GoalForwardPointsFactor
	}

	return
}

func (p Player) ScoreForAssists() (score int) {
	score += p.Stats.Assists * AssistPointsFactor

	return
}

func (p Player) Score() (score int) {
	score += p.ScoreForGoals()

	score += p.ScoreForAssists()

	slog.Debug("player score", "name", p.Name, "score", score, "assists", p.Stats.Assists, "goals", p.Stats.Goals, "position", p.Position)

	return
}

func (c Club) ScoreForWins() (score int) {
	score += c.Stats.Wins * WinPointsFactor

	return
}

func (c Club) ScoreForLossesInOT() (score int) {
	score += c.Stats.LossesInOT * LossesInOTPointsFactor

	return
}

func (c Club) Score() (score int) {
	score += c.ScoreForWins()

	score += c.ScoreForLossesInOT()

	slog.Debug("club score", "name", c.Abbrev, "wins", c.Stats.Wins, "lossesInOT", c.Stats.LossesInOT, "inc", c.Stats.Wins*WinPointsFactor+c.Stats.LossesInOT*LossesInOTPointsFactor)

	return
}

func (t Team) ScoreForGoals() (score int) {
	for _, p := range t.Players {
		score += p.ScoreForGoals()
	}

	// Special Quebec Round
	v, err := t.QuebecRound.ScoreForQuebecRound()
	if err != nil {
		slog.Warn("quebec round failed", "err", err)
	}

	score += v

	return
}

func (t Team) ScoreForAssists() (score int) {
	for _, p := range t.Players {
		score += p.ScoreForAssists()
	}

	return
}

func (t Team) ScoreForWins() (score int) {
	for _, c := range t.Clubs {
		score += c.ScoreForWins()
	}

	return
}

func (t Team) ScoreForLossesInOT() (score int) {
	for _, c := range t.Clubs {
		score += c.ScoreForLossesInOT()
	}

	return
}
