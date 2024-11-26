package teams

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func inputFromDisk(fname string) []byte {
	f, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	return f
}

func inputFromRemote(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return data
}

func FromFile(s string) []Team {
	var data []byte

	if strings.HasPrefix(s, "http") {
		data = inputFromRemote(s)
	} else {
		data = inputFromDisk(s)
	}

	var inputTeams []InputTeam

	err := json.Unmarshal(data, &inputTeams)
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
			t.Players[i].Name = p
		}

		for i, c := range inputTeam.Clubs {
			t.Clubs[i].Abbrev = c
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
