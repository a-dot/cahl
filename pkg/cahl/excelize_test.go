package cahl

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func TestExcelize(t *testing.T) {
	tmpDir := t.TempDir()
	excelPath := filepath.Join(tmpDir, "test_output.xlsx")

	makeTestTeam := func(name, manager, qcPlayer string, qcGoals int) Team {
		players := make([]*Player, TeamPlayerCount)
		for i := 0; i < TeamPlayerCount; i++ {
			players[i] = &Player{
				Name:     fmt.Sprintf("Player %d of %s", i+1, name),
				Position: Forward,
				Stats: PlayerStats{
					Goals:   2, // 2 * 2 = 4
					Assists: 1, // 1 * 1 = 1
				}, // total per player = 5 points
			}
		}

		clubs := make([]*Club, TeamClubCount)
		for i := 0; i < TeamClubCount; i++ {
			clubs[i] = &Club{
				Name:   fmt.Sprintf("Club %d", i+1),
				Abbrev: fmt.Sprintf("C%d", i+1),
				Stats: ClubStats{
					Wins:       3, // 3 * 2 = 6
					LossesInOT: 1, // 1 * 1 = 1
				}, // total per club = 7 points
			}
		}

		return Team{
			Name:    name,
			Manager: manager,
			Players: players,
			Clubs:   clubs,
			QuebecRound: Player{
				Name:  qcPlayer,
				Stats: PlayerStats{Goals: qcGoals},
			},
		}
	}

	// Team 1: 9 players * 5 = 45. 3 clubs * 7 = 21.
	// Claude Giroux threshold 14 -> 14 * 1 * 3 = 42.
	// Total: 45 + 21 + 42 = 108.
	// Goals: 9 * (2 * 2) = 36. Wins: 3 * (3 * 2) = 18. QC: 42. B/V = 36 + 18 + 42 = 96.
	// Assists: 9 * 1 = 9.
	// OT losses: 3 * 1 = 3.
	// Total: 96 + 9 + 3 = 108.
	t1 := makeTestTeam("Les Champions", "Manager A", "Claude Giroux", 14)

	// Team 2: Claude Giroux below threshold (10 goals) -> 10 * 1 = 10.
	// Total: 45 + 21 + 10 = 76.
	// B/V: 36 + 18 + 10 = 64. Assists: 9. OT losses: 3. Total = 76.
	t2 := makeTestTeam("Les Coureurs", "Manager B", "Claude Giroux", 10)

	// Team 3: Alexis Lafrenière threshold 27 (27 goals) -> 27 * 3 = 81.
	// Total: 45 + 21 + 81 = 147.
	t3 := makeTestTeam("Les Étoiles", "Manager C", "Alexis Lafrenière", 27)

	teams := []Team{t1, t2, t3}
	ranking := CreateRanking(teams)
	prevRanking := Ranking{}

	Excelize(teams, ranking, prevRanking, excelPath)

	_, err := os.Stat(excelPath)
	require.NoError(t, err)

	f, err := excelize.OpenFile(excelPath)
	require.NoError(t, err)
	defer f.Close()

	// Verify sheets exist
	sheetList := f.GetSheetList()
	require.Contains(t, sheetList, "Classement")
	require.Contains(t, sheetList, "Equipes")

	// --- Check "Classement" sheet ---
	// Rank 1 should be Team 3 (Score 147)
	// Rank 2 should be Team 1 (Score 108)
	// Rank 3 should be Team 2 (Score 76)
	scoreF2, err := f.GetCellValue("Classement", "F2")
	require.NoError(t, err)
	require.Equal(t, "147", scoreF2)

	scoreF3, err := f.GetCellValue("Classement", "F3")
	require.NoError(t, err)
	require.Equal(t, "108", scoreF3)

	scoreF4, err := f.GetCellValue("Classement", "F4")
	require.NoError(t, err)
	require.Equal(t, "76", scoreF4)

	// Check that Column C + Column D + Column E == Column F for each row
	for row := 2; row <= 4; row++ {
		cVal, err := f.GetCellValue("Classement", fmt.Sprintf("C%d", row))
		require.NoError(t, err)
		dVal, err := f.GetCellValue("Classement", fmt.Sprintf("D%d", row))
		require.NoError(t, err)
		eVal, err := f.GetCellValue("Classement", fmt.Sprintf("E%d", row))
		require.NoError(t, err)
		fVal, err := f.GetCellValue("Classement", fmt.Sprintf("F%d", row))
		require.NoError(t, err)

		c, err := strconv.Atoi(cVal)
		require.NoError(t, err)
		d, err := strconv.Atoi(dVal)
		require.NoError(t, err)
		e, err := strconv.Atoi(eVal)
		require.NoError(t, err)
		fScore, err := strconv.Atoi(fVal)
		require.NoError(t, err)

		require.Equal(t, fScore, c+d+e, "Row %d: C (%d) + D (%d) + E (%d) must equal F (%d)", row, c, d, e, fScore)
	}

	// --- Check "Equipes" sheet ---
	// Team 0 (t1): starts at A1
	teamNameA1, err := f.GetCellValue("Equipes", "A1")
	require.NoError(t, err)
	require.Equal(t, "Les Champions (Manager A)", teamNameA1)

	// Quebec round player for Team 0 is at row 17
	qcA17, err := f.GetCellValue("Equipes", "A17")
	require.NoError(t, err)
	require.Equal(t, "Claude Giroux", qcA17)

	// Quebec round goals points at B17
	qcB17, err := f.GetCellValue("Equipes", "B17")
	require.NoError(t, err)
	require.Equal(t, "42", qcB17)

	// Quebec round total points at E17
	qcE17, err := f.GetCellValue("Equipes", "E17")
	require.NoError(t, err)
	require.Equal(t, "42", qcE17)

	// Team 1 (t2): starts at G1
	teamNameG1, err := f.GetCellValue("Equipes", "G1")
	require.NoError(t, err)
	require.Equal(t, "Les Coureurs (Manager B)", teamNameG1)

	// Quebec round player for Team 1 is at row 17, col G
	qcG17, err := f.GetCellValue("Equipes", "G17")
	require.NoError(t, err)
	require.Equal(t, "Claude Giroux", qcG17)

	qcH17, err := f.GetCellValue("Equipes", "H17")
	require.NoError(t, err)
	require.Equal(t, "10", qcH17)

	qcK17, err := f.GetCellValue("Equipes", "K17")
	require.NoError(t, err)
	require.Equal(t, "10", qcK17)

	// Team 2 (t3): starts at row 21 (due to rowOffset = 18), col A
	teamNameA21, err := f.GetCellValue("Equipes", "A21")
	require.NoError(t, err)
	require.Equal(t, "Les Étoiles (Manager C)", teamNameA21)

	// Quebec round player for Team 2 is at row 37 (21 + 16)
	qcA37, err := f.GetCellValue("Equipes", "A37")
	require.NoError(t, err)
	require.Equal(t, "Alexis Lafrenière", qcA37)

	qcB37, err := f.GetCellValue("Equipes", "B37")
	require.NoError(t, err)
	require.Equal(t, "81", qcB37)

	qcE37, err := f.GetCellValue("Equipes", "E37")
	require.NoError(t, err)
	require.Equal(t, "81", qcE37)
}

