package teams

import (
	"encoding/json"
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

		ret = append(ret, t)
	}

	return ret
}
