package cahl

import (
	"slices"
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

func TestRankSort(t *testing.T) {
	a := Rank{
		Team: Team{
			Name: "A",
			Players: []*Player{
				{Position: Forward, Stats: PlayerStats{Goals: 2}},
			},
		},
	}
	a.Score = a.Team.Score()

	b := Rank{
		Team: Team{
			Name: "B",
			Players: []*Player{
				{Position: Forward, Stats: PlayerStats{Goals: 3}},
			},
		},
	}
	b.Score = b.Team.Score()

	teams := []Rank{a, b}

	slices.SortFunc(teams, rankSort)

	require.Equal(t, "B", teams[0].Team.Name)
}

func TestRankSortTieBreaker(t *testing.T) {
	a := Rank{
		Team: Team{
			Name: "A",
			Players: []*Player{
				{Position: Forward, Stats: PlayerStats{Goals: 2, Assists: 2}},
			},
		},
	}
	a.Score = a.Team.Score()

	b := Rank{
		Team: Team{
			Name: "B",
			Players: []*Player{
				{Position: Forward, Stats: PlayerStats{Goals: 3}},
			},
		},
	}
	b.Score = b.Team.Score()

	teams := []Rank{a, b}

	slices.SortFunc(teams, rankSort)

	require.Equal(t, "B", teams[0].Team.Name)
}
