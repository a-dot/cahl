package cahl

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadTeamsFromTestFile(t *testing.T) {
	testFilePath := filepath.Join("..", "..", "teams_test.json")
	if _, err := os.Stat(testFilePath); errors.Is(err, os.ErrNotExist) {
		t.Skipf("skipping test: %s does not exist", testFilePath)
	}

	teams, err := LoadTeams(testFilePath)
	require.NoError(t, err)
	require.Len(t, teams, 15)

	for _, team := range teams {
		require.NoError(t, team.Valid(), "team %s should be valid", team.Name)
		require.NotEmpty(t, team.Name)
		require.NotEmpty(t, team.Manager)
		require.Len(t, team.Players, TeamPlayerCount)
		require.Len(t, team.Clubs, TeamClubCount)
		require.NotEmpty(t, team.QuebecRound.Name)

		// Verify each quebec_round player has a threshold defined
		_, err := team.QuebecRound.ScoreForQuebecRound()
		require.NoError(t, err, "quebec_round player %s should be in threshold map", team.QuebecRound.Name)
	}
}
