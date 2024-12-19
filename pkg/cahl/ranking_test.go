package cahl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPositionUnordered(t *testing.T) {
	r := Rank{
		Team: Team{
			Name: "test-01",
		},
		Score: 10,
	}

	ranking := Ranking{
		Teams: []Rank{
			r,
			{Team: Team{Name: "test-02"}, Score: 20},
		},
	}

	require.Equal(t, 2, ranking.Position(r))
}

func TestPositionOrdered(t *testing.T) {
	r := Rank{
		Team: Team{
			Name: "test-01",
		},
		Score: 10,
	}

	ranking := Ranking{
		Teams: []Rank{
			{Team: Team{Name: "test-02"}, Score: 20},
			r,
		},
	}

	require.Equal(t, 2, ranking.Position(r))
}

func TestDeltaFrom(t *testing.T) {
	ra := Rank{
		Team: Team{
			Name: "test-01",
		},
	}

	rankingA := Ranking{
		Teams: []Rank{
			{Team: Team{Name: "team-one"}, Score: 20},
			{Team: Team{Name: "test-01"}, Score: 10},
		},
	}

	rankingB := Ranking{
		Teams: []Rank{
			{Team: Team{Name: "test-01"}, Score: 21},
			{Team: Team{Name: "team-one"}, Score: 20},
		},
	}

	res := ra.DeltaFrom(rankingB, rankingA)

	require.Equal(t, 1, res.Position)
	require.Equal(t, 11, res.Score)
}

func TestDeltaFromNegativePosition(t *testing.T) {
	ra := Rank{
		Team: Team{
			Name: "test-01",
		},
	}

	rankingA := Ranking{
		Teams: []Rank{
			{Team: Team{Name: "test-01"}, Score: 21},
			{Team: Team{Name: "team-one"}, Score: 20},
		},
	}

	rankingB := Ranking{
		Teams: []Rank{
			{Team: Team{Name: "team-one"}, Score: 45},
			{Team: Team{Name: "test-01"}, Score: 25},
		},
	}

	res := ra.DeltaFrom(rankingB, rankingA)

	require.Equal(t, -1, res.Position)
	require.Equal(t, 4, res.Score)
}