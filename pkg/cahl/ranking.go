package cahl

import (
	"fmt"
	"log/slog"
	"slices"
)

type Rank struct {
	Team  Team
	Score int
	Delta DeltaFromPrev
}

type DeltaFromPrev struct {
	Score    int
	Position int
}

func (r Rank) String() string {
	// TODO include Delta
	return fmt.Sprintf("(%d)%s", r.Score, r.Team.Name)
}

func CreateRanking(teams []Team) (ranking []Rank) {
	ranking = make([]Rank, 0, len(teams))

	for _, team := range teams {
		ranking = append(ranking, Rank{
			Team:  team,
			Score: team.Score(),
		})
	}

	slices.SortFunc(ranking, func(a, b Rank) int {
		return b.Score - a.Score
	})

	slog.Debug("ranking", "ranking", ranking)

	//TODO Deltas

	return
}
