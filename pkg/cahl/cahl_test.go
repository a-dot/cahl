package cahl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScoreForQuebecRound(t *testing.T) {
	tests := []struct {
		name          string
		player        Player
		expectedScore int
		expectedError string
	}{
		{
			name: "Unknown player returns error",
			player: Player{
				Name:  "Connor McDavid",
				Stats: PlayerStats{Goals: 30},
			},
			expectedScore: 0,
			expectedError: "invalid name 'Connor McDavid'",
		},
		{
			name: "Empty name returns error",
			player: Player{
				Name:  "",
				Stats: PlayerStats{Goals: 10},
			},
			expectedScore: 0,
			expectedError: "invalid name ''",
		},
		// Claude Giroux (threshold: 14)
		{
			name: "Claude Giroux with 0 goals",
			player: Player{
				Name:  "Claude Giroux",
				Stats: PlayerStats{Goals: 0},
			},
			expectedScore: 0,
		},
		{
			name: "Claude Giroux below threshold (13 goals)",
			player: Player{
				Name:  "Claude Giroux",
				Stats: PlayerStats{Goals: 13},
			},
			expectedScore: 13,
		},
		{
			name: "Claude Giroux at threshold (14 goals)",
			player: Player{
				Name:  "Claude Giroux",
				Stats: PlayerStats{Goals: 14},
			},
			expectedScore: 42, // 14 * 1 * 3
		},
		{
			name: "Claude Giroux above threshold (15 goals)",
			player: Player{
				Name:  "Claude Giroux",
				Stats: PlayerStats{Goals: 15},
			},
			expectedScore: 45, // 15 * 1 * 3
		},
		// Philippe Danault (threshold: 14)
		{
			name: "Philippe Danault below threshold (10 goals)",
			player: Player{
				Name:  "Philippe Danault",
				Stats: PlayerStats{Goals: 10},
			},
			expectedScore: 10,
		},
		{
			name: "Philippe Danault at threshold (14 goals)",
			player: Player{
				Name:  "Philippe Danault",
				Stats: PlayerStats{Goals: 14},
			},
			expectedScore: 42,
		},
		// Zachary Bolduc (threshold: 14)
		{
			name: "Zachary Bolduc below threshold (13 goals)",
			player: Player{
				Name:  "Zachary Bolduc",
				Stats: PlayerStats{Goals: 13},
			},
			expectedScore: 13,
		},
		{
			name: "Zachary Bolduc at threshold (14 goals)",
			player: Player{
				Name:  "Zachary Bolduc",
				Stats: PlayerStats{Goals: 14},
			},
			expectedScore: 42,
		},
		// Alexis Lafrenière (threshold: 27)
		{
			name: "Alexis Lafrenière below threshold (26 goals)",
			player: Player{
				Name:  "Alexis Lafrenière",
				Stats: PlayerStats{Goals: 26},
			},
			expectedScore: 26,
		},
		{
			name: "Alexis Lafrenière at threshold (27 goals)",
			player: Player{
				Name:  "Alexis Lafrenière",
				Stats: PlayerStats{Goals: 27},
			},
			expectedScore: 81, // 27 * 1 * 3
		},
		{
			name: "Alexis Lafrenière above threshold (30 goals)",
			player: Player{
				Name:  "Alexis Lafrenière",
				Stats: PlayerStats{Goals: 30},
			},
			expectedScore: 90, // 30 * 1 * 3
		},
		// Pierre-Luc Dubois (threshold: 21)
		{
			name: "Pierre-Luc Dubois below threshold (20 goals)",
			player: Player{
				Name:  "Pierre-Luc Dubois",
				Stats: PlayerStats{Goals: 20},
			},
			expectedScore: 20,
		},
		{
			name: "Pierre-Luc Dubois at threshold (21 goals)",
			player: Player{
				Name:  "Pierre-Luc Dubois",
				Stats: PlayerStats{Goals: 21},
			},
			expectedScore: 63, // 21 * 1 * 3
		},
		{
			name: "Pierre-Luc Dubois above threshold (25 goals)",
			player: Player{
				Name:  "Pierre-Luc Dubois",
				Stats: PlayerStats{Goals: 25},
			},
			expectedScore: 75, // 25 * 1 * 3
		},
		// Jonathan Huberdeau (threshold: 23)
		{
			name: "Jonathan Huberdeau below threshold (22 goals)",
			player: Player{
				Name:  "Jonathan Huberdeau",
				Stats: PlayerStats{Goals: 22},
			},
			expectedScore: 22,
		},
		{
			name: "Jonathan Huberdeau at threshold (23 goals)",
			player: Player{
				Name:  "Jonathan Huberdeau",
				Stats: PlayerStats{Goals: 23},
			},
			expectedScore: 69, // 23 * 1 * 3
		},
		{
			name: "Jonathan Huberdeau above threshold (24 goals)",
			player: Player{
				Name:  "Jonathan Huberdeau",
				Stats: PlayerStats{Goals: 24},
			},
			expectedScore: 72, // 24 * 1 * 3
		},
		// Assists should not affect Quebec round score
		{
			name: "Assists do not affect Quebec round score",
			player: Player{
				Name: "Claude Giroux",
				Stats: PlayerStats{
					Goals:   14,
					Assists: 20,
				},
			},
			expectedScore: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, err := tt.player.ScoreForQuebecRound()

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				require.Equal(t, 0, score)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedScore, score)
			}
		})
	}
}

func TestTeamScoreForGoals(t *testing.T) {
	tests := []struct {
		name          string
		team          Team
		expectedScore int
	}{
		{
			name:          "Empty team with no players or Quebec round",
			team:          Team{},
			expectedScore: 0,
		},
		{
			name: "Forwards only (2 points per goal)",
			team: Team{
				Players: []*Player{
					{Name: "Forward 1", Position: Forward, Stats: PlayerStats{Goals: 5}},
					{Name: "Forward 2", Position: Forward, Stats: PlayerStats{Goals: 3}},
				},
			},
			expectedScore: 16, // (5 + 3) * 2
		},
		{
			name: "Defencemen only (3 points per goal)",
			team: Team{
				Players: []*Player{
					{Name: "Defenceman 1", Position: Defence, Stats: PlayerStats{Goals: 2}},
					{Name: "Defenceman 2", Position: Defence, Stats: PlayerStats{Goals: 4}},
				},
			},
			expectedScore: 18, // (2 + 4) * 3
		},
		{
			name: "Forwards and Defencemen mixed",
			team: Team{
				Players: []*Player{
					{Name: "Forward 1", Position: Forward, Stats: PlayerStats{Goals: 4}},    // 8
					{Name: "Defenceman 1", Position: Defence, Stats: PlayerStats{Goals: 2}}, // 6
				},
			},
			expectedScore: 14, // 8 + 6
		},
		{
			name: "Team with Quebec round below threshold",
			team: Team{
				Players: []*Player{
					{Name: "Forward 1", Position: Forward, Stats: PlayerStats{Goals: 2}}, // 4
				},
				QuebecRound: Player{
					Name:  "Claude Giroux",
					Stats: PlayerStats{Goals: 10}, // threshold 14 -> 10 * 1 = 10
				},
			},
			expectedScore: 14, // 4 + 10
		},
		{
			name: "Team with Quebec round at or above threshold",
			team: Team{
				Players: []*Player{
					{Name: "Forward 1", Position: Forward, Stats: PlayerStats{Goals: 2}}, // 4
				},
				QuebecRound: Player{
					Name:  "Claude Giroux",
					Stats: PlayerStats{Goals: 14}, // threshold 14 -> 14 * 1 * 3 = 42
				},
			},
			expectedScore: 46, // 4 + 42
		},
		{
			name: "Team with invalid Quebec round player handles error and adds 0",
			team: Team{
				Players: []*Player{
					{Name: "Forward 1", Position: Forward, Stats: PlayerStats{Goals: 3}}, // 6
				},
				QuebecRound: Player{
					Name:  "Unknown Player",
					Stats: PlayerStats{Goals: 50},
				},
			},
			expectedScore: 6, // 6 + 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.team.ScoreForGoals()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestTeamScoreForAssists(t *testing.T) {
	tests := []struct {
		name          string
		team          Team
		expectedScore int
	}{
		{
			name:          "Empty team with no players",
			team:          Team{},
			expectedScore: 0,
		},
		{
			name: "Players with zero assists",
			team: Team{
				Players: []*Player{
					{Name: "Player 1", Position: Forward, Stats: PlayerStats{Assists: 0}},
					{Name: "Player 2", Position: Defence, Stats: PlayerStats{Assists: 0}},
				},
			},
			expectedScore: 0,
		},
		{
			name: "Multiple players with assists (1 point per assist)",
			team: Team{
				Players: []*Player{
					{Name: "Forward 1", Position: Forward, Stats: PlayerStats{Assists: 5}},
					{Name: "Forward 2", Position: Forward, Stats: PlayerStats{Assists: 3}},
					{Name: "Defence 1", Position: Defence, Stats: PlayerStats{Assists: 10}},
				},
			},
			expectedScore: 18, // 5 + 3 + 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.team.ScoreForAssists()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestTeamScoreForWins(t *testing.T) {
	tests := []struct {
		name          string
		team          Team
		expectedScore int
	}{
		{
			name:          "Empty team with no clubs",
			team:          Team{},
			expectedScore: 0,
		},
		{
			name: "Clubs with zero wins",
			team: Team{
				Clubs: []*Club{
					{Name: "Team A", Abbrev: "TMA", Stats: ClubStats{Wins: 0}},
					{Name: "Team B", Abbrev: "TMB", Stats: ClubStats{Wins: 0}},
				},
			},
			expectedScore: 0,
		},
		{
			name: "Clubs with wins (2 points per win)",
			team: Team{
				Clubs: []*Club{
					{Name: "Montreal Canadiens", Abbrev: "MTL", Stats: ClubStats{Wins: 10}},
					{Name: "Toronto Maple Leafs", Abbrev: "TOR", Stats: ClubStats{Wins: 20}},
					{Name: "Boston Bruins", Abbrev: "BOS", Stats: ClubStats{Wins: 5}},
				},
			},
			expectedScore: 70, // (10 + 20 + 5) * 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.team.ScoreForWins()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestTeamScoreForLossesInOT(t *testing.T) {
	tests := []struct {
		name          string
		team          Team
		expectedScore int
	}{
		{
			name:          "Empty team with no clubs",
			team:          Team{},
			expectedScore: 0,
		},
		{
			name: "Clubs with zero OT losses",
			team: Team{
				Clubs: []*Club{
					{Name: "Team A", Abbrev: "TMA", Stats: ClubStats{LossesInOT: 0}},
					{Name: "Team B", Abbrev: "TMB", Stats: ClubStats{LossesInOT: 0}},
				},
			},
			expectedScore: 0,
		},
		{
			name: "Clubs with OT losses (1 point per OT loss)",
			team: Team{
				Clubs: []*Club{
					{Name: "Montreal Canadiens", Abbrev: "MTL", Stats: ClubStats{LossesInOT: 3}},
					{Name: "Toronto Maple Leafs", Abbrev: "TOR", Stats: ClubStats{LossesInOT: 5}},
					{Name: "Boston Bruins", Abbrev: "BOS", Stats: ClubStats{LossesInOT: 2}},
				},
			},
			expectedScore: 10, // (3 + 5 + 2) * 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.team.ScoreForLossesInOT()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestTeamValid(t *testing.T) {
	makePlayers := func(n int) []*Player {
		p := make([]*Player, n)
		for i := range p {
			p[i] = &Player{Name: fmt.Sprintf("Player %d", i+1)}
		}
		return p
	}

	makeClubs := func(n int) []*Club {
		c := make([]*Club, n)
		for i := range c {
			c[i] = &Club{Name: fmt.Sprintf("Club %d", i+1)}
		}
		return c
	}

	tests := []struct {
		name          string
		team          Team
		expectedError string
	}{
		{
			name: fmt.Sprintf("Valid team with %d players and %d clubs", TeamPlayerCount, TeamClubCount),
			team: Team{
				Name:    "Champs",
				Players: makePlayers(TeamPlayerCount),
				Clubs:   makeClubs(TeamClubCount),
			},
		},
		{
			name: fmt.Sprintf("Too few players (%d players)", TeamPlayerCount-1),
			team: Team{
				Name:    "Understaffed",
				Players: makePlayers(TeamPlayerCount - 1),
				Clubs:   makeClubs(TeamClubCount),
			},
			expectedError: fmt.Sprintf("team 'Understaffed' has the wrong number of players (%d)", TeamPlayerCount-1),
		},
		{
			name: fmt.Sprintf("Too many players (%d players)", TeamPlayerCount+1),
			team: Team{
				Name:    "Overstaffed",
				Players: makePlayers(TeamPlayerCount + 1),
				Clubs:   makeClubs(TeamClubCount),
			},
			expectedError: fmt.Sprintf("team 'Overstaffed' has the wrong number of players (%d)", TeamPlayerCount+1),
		},
		{
			name: fmt.Sprintf("Too few clubs (%d clubs)", TeamClubCount-1),
			team: Team{
				Name:    "Missing Club",
				Players: makePlayers(TeamPlayerCount),
				Clubs:   makeClubs(TeamClubCount - 1),
			},
			expectedError: fmt.Sprintf("team 'Missing Club' has the wrong number of clubs (%d)", TeamClubCount-1),
		},
		{
			name: fmt.Sprintf("Too many clubs (%d clubs)", TeamClubCount+1),
			team: Team{
				Name:    "Extra Club",
				Players: makePlayers(TeamPlayerCount),
				Clubs:   makeClubs(TeamClubCount + 1),
			},
			expectedError: fmt.Sprintf("team 'Extra Club' has the wrong number of clubs (%d)", TeamClubCount+1),
		},
		{
			name: "Both wrong player and club count fails on players check first",
			team: Team{
				Name:    "Empty",
				Players: nil,
				Clubs:   nil,
			},
			expectedError: "team 'Empty' has the wrong number of players (0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.team.Valid()
			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTeamScore(t *testing.T) {
	tests := []struct {
		name          string
		team          Team
		expectedScore int
	}{
		{
			name:          "Empty team",
			team:          Team{Name: "Empty"},
			expectedScore: 0,
		},
		{
			name: "Team with players and clubs",
			team: Team{
				Name: "Full Team",
				Players: []*Player{
					// Forward: 5 goals * 2 = 10, 3 assists * 1 = 3 -> 13
					{Name: "Forward", Position: Forward, Stats: PlayerStats{Goals: 5, Assists: 3}},
					// Defence: 2 goals * 3 = 6, 4 assists * 1 = 4 -> 10
					{Name: "Defence", Position: Defence, Stats: PlayerStats{Goals: 2, Assists: 4}},
				},
				Clubs: []*Club{
					// Club: 10 wins * 2 = 20, 3 OT losses * 1 = 3 -> 23
					{Name: "Club 1", Stats: ClubStats{Wins: 10, LossesInOT: 3}},
				},
			},
			// 13 + 10 + 23 = 46
			expectedScore: 46,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.team.Score()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestPlayerScoreForGoals(t *testing.T) {
	tests := []struct {
		name          string
		player        Player
		expectedScore int
	}{
		{
			name: "Forward with 0 goals",
			player: Player{
				Position: Forward,
				Stats:    PlayerStats{Goals: 0},
			},
			expectedScore: 0,
		},
		{
			name: "Forward with goals (2 points each)",
			player: Player{
				Position: Forward,
				Stats:    PlayerStats{Goals: 5},
			},
			expectedScore: 10,
		},
		{
			name: "Defence with 0 goals",
			player: Player{
				Position: Defence,
				Stats:    PlayerStats{Goals: 0},
			},
			expectedScore: 0,
		},
		{
			name: "Defence with goals (3 points each)",
			player: Player{
				Position: Defence,
				Stats:    PlayerStats{Goals: 4},
			},
			expectedScore: 12,
		},
		{
			name: "Unknown position defaults to forward scoring factor (2 points each)",
			player: Player{
				Position: Unknown,
				Stats:    PlayerStats{Goals: 3},
			},
			expectedScore: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.player.ScoreForGoals()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestPlayerScoreForAssists(t *testing.T) {
	tests := []struct {
		name          string
		player        Player
		expectedScore int
	}{
		{
			name: "Player with 0 assists",
			player: Player{
				Stats: PlayerStats{Assists: 0},
			},
			expectedScore: 0,
		},
		{
			name: "Player with assists (1 point each)",
			player: Player{
				Stats: PlayerStats{Assists: 7},
			},
			expectedScore: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.player.ScoreForAssists()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestPlayerScore(t *testing.T) {
	tests := []struct {
		name          string
		player        Player
		expectedScore int
	}{
		{
			name: "Forward goals and assists",
			player: Player{
				Name:     "Forward Star",
				Position: Forward,
				Stats:    PlayerStats{Goals: 4, Assists: 6},
			},
			expectedScore: 14, // (4 * 2) + (6 * 1)
		},
		{
			name: "Defence goals and assists",
			player: Player{
				Name:     "Top D-man",
				Position: Defence,
				Stats:    PlayerStats{Goals: 3, Assists: 5},
			},
			expectedScore: 14, // (3 * 3) + (5 * 1)
		},
		{
			name: "Player with zero stats",
			player: Player{
				Name:     "Rookie",
				Position: Forward,
				Stats:    PlayerStats{Goals: 0, Assists: 0},
			},
			expectedScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.player.Score()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestClubScoreForWins(t *testing.T) {
	tests := []struct {
		name          string
		club          Club
		expectedScore int
	}{
		{
			name:          "Club with 0 wins",
			club:          Club{Stats: ClubStats{Wins: 0}},
			expectedScore: 0,
		},
		{
			name:          "Club with wins (2 points per win)",
			club:          Club{Stats: ClubStats{Wins: 12}},
			expectedScore: 24, // 12 * 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.club.ScoreForWins()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestClubScoreForLossesInOT(t *testing.T) {
	tests := []struct {
		name          string
		club          Club
		expectedScore int
	}{
		{
			name:          "Club with 0 OT losses",
			club:          Club{Stats: ClubStats{LossesInOT: 0}},
			expectedScore: 0,
		},
		{
			name:          "Club with OT losses (1 point per OT loss)",
			club:          Club{Stats: ClubStats{LossesInOT: 7}},
			expectedScore: 7, // 7 * 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.club.ScoreForLossesInOT()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestClubScore(t *testing.T) {
	tests := []struct {
		name          string
		club          Club
		expectedScore int
	}{
		{
			name:          "Club with 0 wins and 0 OT losses",
			club:          Club{Stats: ClubStats{Wins: 0, LossesInOT: 0}},
			expectedScore: 0,
		},
		{
			name:          "Club with wins and OT losses",
			club:          Club{Stats: ClubStats{Wins: 8, LossesInOT: 3}},
			expectedScore: 19, // (8 * 2) + (3 * 1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := tt.club.Score()
			require.Equal(t, tt.expectedScore, score)
		})
	}
}

