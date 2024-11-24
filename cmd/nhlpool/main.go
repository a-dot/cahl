package main

import (
	"fmt"
	"log/slog"
	"nhlpool/pkg/teams"
	"os"
	"slices"
)

var CURRENTSEASON = "20242025" //TODO make command line argument

type Ranking struct {
	teamName string
	score    int
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	teams := teams.FromFile("../../resources/teams.json")

	ranking := make([]Ranking, 0, len(teams))

	for _, team := range teams {
		team.PopulatePlayersStats(CURRENTSEASON)

		team.PopulateClubsStats(CURRENTSEASON)

		score := team.Score()

		ranking = append(ranking, Ranking{team.Name, score})

		slog.Debug("total score for team", "team", team.Name, "score", score)
	}

	slices.SortFunc(ranking, func(a, b Ranking) int {
		return b.score - a.score
	})

	fmt.Println(ranking)
}
