package cahl

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"
)

type Rank struct {
	Team  Team `json:"team"`
	Score int  `json:"score"`
}

type DeltaFromPrev struct {
	Score    int
	Position int
}

type Ranking struct {
	Timestamp time.Time `json:"timestamp"`
	Teams     []Rank    `json:"teams"`
}

func (r Rank) String() string {
	return fmt.Sprintf("(%d)%s", r.Score, r.Team.Name)
}

func CreateRanking(teams []Team) Ranking {
	ranking := make([]Rank, 0, len(teams))

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

	return Ranking{
		Timestamp: time.Now(),
		Teams:     ranking,
	}
}

func (t Rank) DeltaFrom(cur, prev Ranking) DeltaFromPrev {
	curPosition := cur.Position(t)
	prevPosition := prev.Position(t)

	curScore, err := cur.TeamScore(t)
	if err != nil {
		panic(err)
	}

	prevScore, err := prev.TeamScore(t)
	if err != nil {
		panic(err)
	}

	return DeltaFromPrev{
		Score:    curScore - prevScore,
		Position: prevPosition - curPosition,
	}
}

func (ranking Ranking) Position(r Rank) int {
	slices.SortFunc(ranking.Teams, func(a, b Rank) int {
		return b.Score - a.Score
	})

	for i, t := range ranking.Teams {
		if t.Team.Name == r.Team.Name {
			return i + 1
		}
	}

	return 0
}

func (ranking Ranking) TeamScore(r Rank) (int, error) {
	for _, rank := range ranking.Teams {
		if rank.Team.Name == r.Team.Name {
			return rank.Score, nil
		}
	}

	return 0, errors.New("not found")
}
