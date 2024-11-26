package teams

import (
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

type Player struct {
	name     string
	id       uint64
	position Position

	goals   int
	assists int
}

type Club struct {
	name   string
	abbrev string

	wins       int
	lossesInOT int
}

type Team struct {
	Name    string
	Players []Player
	Clubs   []Club
}

func (t Team) Score() int {
	var score int

	for _, p := range t.Players {
		if p.position == Defence {
			slog.Debug("defence goals score update", "name", p.name, "goals", p.goals, "inc", p.goals*3)
			score += p.goals * 3
		} else {
			slog.Debug("forward goals score update", "name", p.name, "goals", p.goals, "inc", p.goals*2)
			score += p.goals * 2
		}

		slog.Debug("assists score update", "name", p.name, "assists", p.assists, "inc", p.assists)
		score += p.assists
	}

	for _, c := range t.Clubs {
		slog.Debug("team score inc", "name", c.abbrev, "wins", c.wins, "lossesInOT", c.lossesInOT, "inc", c.wins*2+c.lossesInOT)
		score += c.wins*2 + c.lossesInOT
	}

	slog.Debug("calculated team score", "team", t.Name, "score", score)

	return score
}

func (t *Team) PopulatePlayersStats(season string) {
	for i, p := range t.Players {
		// Search for player
		id := api.SearchPlayer(p.name, season)

		if id == 0 {
			slog.Error("unable to find player", "team", t.Name, "player name", p.name)
			panic("unable to find player")
		}

		// Get player stats for current season
		playerStats := api.PlayerStats(id, season)
		if playerStats == nil {
			slog.Error("unable to find stats for player", "team", t.Name, "player name", p.name, "id", id)
			panic("unable to find stats for player")
		}

		switch playerStats.Position {
		case "D":
			t.Players[i].position = Defence
		default:
			t.Players[i].position = Forward
		}

		t.Players[i].goals = playerStats.Goals
		t.Players[i].assists = playerStats.Assists
	}
}

func (t *Team) PopulateClubsStats(season string) {
	// Get stats for all teams
	teamStats := api.StatsForAllTeams()

	for i, c := range t.Clubs {
		C, found := teamStats[c.abbrev]
		if !found {
			slog.Error("team not found", "team abbrev", c.abbrev)
			panic("team not found")
		}

		t.Clubs[i].wins = C.Wins
		t.Clubs[i].lossesInOT = C.LossesOT
	}
}
