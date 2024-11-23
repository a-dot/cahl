package main

import (
	"log/slog"
	"nhlpool/pkg/teams"
	"os"
)

var CURRENTSEASON = "20242025" //TODO make command line argument

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	teams := teams.FromFile("../../resources/teams.json")

	for _, team := range teams {
		team.PopulatePlayersStats(CURRENTSEASON)

		team.PopulateClubsStats(CURRENTSEASON)

		score := team.Score()

		slog.Debug("total score for team", "team", team.Name, "score", score)
	}
}
