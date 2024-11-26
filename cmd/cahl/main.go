package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/jessevdk/go-flags"

	"cahl/pkg/teams"
)

type Ranking struct {
	teamName string
	score    int
}

var opts struct {
	TeamsFile      string `short:"t" description:"teams file" default:"https://raw.githubusercontent.com/a-dot/cahl-teams/refs/heads/main/teams.json"`
	Season         string `short:"s" description:"season (format is YYYYXXXX)" default:"20242025"`
	DataOutputFile string `short:"d" description:"output json file with information used to calculate ranking"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	teams := teams.FromFile(opts.TeamsFile)

	ranking := make([]Ranking, 0, len(teams))

	for _, team := range teams {
		team.PopulatePlayersStats(opts.Season)

		team.PopulateClubsStats(opts.Season)

		score := team.Score()

		ranking = append(ranking, Ranking{team.Name, score})

		slog.Debug("total score for team", "team", team.Name, "score", score)
	}

	slices.SortFunc(ranking, func(a, b Ranking) int {
		return b.score - a.score
	})

	fmt.Println(ranking)

	if len(opts.DataOutputFile) > 0 {
		teamsData, err := json.Marshal(teams)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(opts.DataOutputFile, teamsData, 0644)
		if err != nil {
			panic(err)
		}
	}
}
