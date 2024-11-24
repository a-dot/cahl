package teams

import (
	"encoding/json"
	"log/slog"
	"os"
)

func FromFile(s string) []Team {
	f, err := os.ReadFile(s)
	if err != nil {
		panic(err)
	}

	var inputTeams []InputTeam

	err = json.Unmarshal(f, &inputTeams)
	if err != nil {
		panic(err)
	}

	ret := make([]Team, 0, len(inputTeams))

	for _, inputTeam := range inputTeams {
		t := Team{
			Name:    inputTeam.Name,
			Players: make([]Player, len(inputTeam.Players)),
			Clubs:   make([]Club, len(inputTeam.Clubs)),
		}

		for i, p := range inputTeam.Players {
			t.Players[i].name = p
		}

		for i, c := range inputTeam.Clubs {
			t.Clubs[i].abbrev = c
		}

		// Sanity checks
		if len(t.Players) != 9 {
			slog.Error("team has the wrong number of players", "team", t.Name, "count", len(t.Players))
			os.Exit(1)
		}

		if len(t.Clubs) != 3 {
			slog.Error("team has the wrong number of clubs", "team", t.Name, "count", len(t.Clubs))
			os.Exit(1)
		}

		ret = append(ret, t)
	}

	return ret
}
