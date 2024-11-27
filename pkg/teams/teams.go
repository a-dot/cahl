package teams

import (
	"encoding/json"
	"errors"
	"log/slog"

	"cahl/pkg/api"
)

type InputTeam struct {
	Name    string   `json:"team_name"`
	Manager string   `json:"manager"`
	Players []string `json:"players"`
	Clubs   []string `json:"teams"`
}

type Position int

const (
	Forward Position = iota
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

func ParsePosition(s string) (Position, error) {
	switch s {
	case "forward":
		return Forward, nil
	case "defence":
		return Defence, nil
	default:
		return 0, errors.New("unknown position")
	}
}

func (p Position) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Position) UnmarshalJSON(data []byte) (err error) {
	var pos string

	if err := json.Unmarshal(data, &pos); err != nil {
		return err
	}

	if *p, err = ParsePosition(pos); err != nil {
		return err
	}

	return nil
}

type Player struct {
	Name     string   `json:"name"`
	ID       uint64   `json:"player_id,omitempty"`
	Position Position `json:"position"`

	Goals   int `json:"goals"`
	Assists int `json:"assists"`
}

type Club struct {
	Name   string `json:"name"`
	Abbrev string `json:"abbrev"`

	Wins       int `json:"wins"`
	LossesInOT int `json:"losses_in_ot"`
}

type Team struct {
	Name    string   `json:"team_name"`
	Players []Player `json:"players"`
	Clubs   []Club   `json:"clubs"`
}

func (t Team) Score() int {
	var score int

	for _, p := range t.Players {
		if p.Position == Defence {
			slog.Debug("defence goals score update", "name", p.Name, "goals", p.Goals, "inc", p.Goals*3)
			score += p.Goals * 3
		} else {
			slog.Debug("forward goals score update", "name", p.Name, "goals", p.Goals, "inc", p.Goals*2)
			score += p.Goals * 2
		}

		slog.Debug("assists score update", "name", p.Name, "assists", p.Assists, "inc", p.Assists)
		score += p.Assists
	}

	for _, c := range t.Clubs {
		slog.Debug("team score inc", "name", c.Abbrev, "wins", c.Wins, "lossesInOT", c.LossesInOT, "inc", c.Wins*2+c.LossesInOT)
		score += c.Wins*2 + c.LossesInOT
	}

	slog.Debug("calculated team score", "team", t.Name, "score", score)

	return score
}

func (t *Team) PopulatePlayersStats(season string) {
	for i, p := range t.Players {
		// Search for player
		id := api.SearchPlayer(p.Name, season)

		if id == 0 {
			slog.Error("unable to find player", "team", t.Name, "player name", p.Name)
			panic("unable to find player")
		}

		t.Players[i].ID = id

		// Get player stats for current season
		playerStats := api.PlayerStats(id, season)
		if playerStats == nil {
			slog.Error("unable to find stats for player", "team", t.Name, "player name", p.Name, "id", id)
			panic("unable to find stats for player")
		}

		switch playerStats.Position {
		case "D":
			t.Players[i].Position = Defence
		default:
			t.Players[i].Position = Forward
		}

		t.Players[i].Goals = playerStats.Goals
		t.Players[i].Assists = playerStats.Assists
	}
}

func (t *Team) PopulateClubsStats(season string) {
	// Get stats for all teams
	teamStats := api.StatsForAllTeams()

	for i, c := range t.Clubs {
		C, found := teamStats[c.Abbrev]
		if !found {
			slog.Error("team not found", "team abbrev", c.Abbrev)
			panic("team not found")
		}

		t.Clubs[i].Wins = C.Wins
		t.Clubs[i].LossesInOT = C.LossesOT
	}
}
