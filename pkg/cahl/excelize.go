package cahl

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func Excelize(teams []Team, curRanking, prevRanking Ranking, outputFile string) {
	f := excelize.NewFile()
	defer f.Close()

	index := classementSheet(f, "Classement", curRanking, prevRanking)
	f.SetActiveSheet(index)

	equipesSheet(f, "Equipes", teams)

	f.DeleteSheet("Sheet1")

	// Save spreadsheet by the given path.
	if err := f.SaveAs(outputFile); err != nil {
		panic(err)
	}
}

var colLetters []rune = []rune{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y', 'Z',
}

func cellName(row int, col rune) string {
	return fmt.Sprintf("%c%d", col, row)
}

func equipesSheet(f *excelize.File, sheetName string, teams []Team) int {
	index, err := f.NewSheet(sheetName)
	if err != nil {
		panic(err)
	}

	f.SetColWidth(sheetName, "A", "A", 20)
	f.SetColWidth(sheetName, "G", "G", 20)

	rowOffset := 0
	for i, team := range teams {
		var col rune
		var row int

		if i%2 == 0 {
			if i > 0 {
				rowOffset += 17
			}

			row = i + 1 + rowOffset
			col = 'A'
		} else {
			row = i + rowOffset
			col = 'G'
		}

		genTeam(f, sheetName, row, col, team)
	}

	return index
}

func genTeam(f *excelize.File, sheetName string, row int, colOrigin rune, team Team) {
	cell := cellName(row, colOrigin)

	// Team name
	f.SetCellValue(sheetName, cell, fmt.Sprintf("%s (%s)", team.Name, team.Manager))

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
	})
	if err != nil {
		panic(err)
	}

	f.SetCellStyle(sheetName, cell, cell, style)

	row += 2

	// Headers
	startCol := colOrigin

	f.SetCellValue(sheetName, cellName(row, colOrigin), "JOUEURS")

	col := colLetters[slices.Index(colLetters, colOrigin)+1]
	f.SetCellValue(sheetName, cellName(row, col), "B")

	col = colLetters[slices.Index(colLetters, col)+1]
	f.SetCellValue(sheetName, cellName(row, col), "P")

	col = colLetters[slices.Index(colLetters, col)+1]
	f.SetCellValue(sheetName, cellName(row, col), "V")

	col = colLetters[slices.Index(colLetters, col)+1]
	f.SetCellValue(sheetName, cellName(row, col), "TOT")

	startRow := row

	row += 2

	// Clubs
	for _, club := range team.Clubs {
		col = colOrigin

		f.SetCellValue(sheetName, cellName(row, col), club.Abbrev)

		col = colLetters[slices.Index(colLetters, col)+1]
		f.SetCellValue(sheetName, cellName(row, col), club.ScoreForWins())

		col = colLetters[slices.Index(colLetters, col)+2]
		f.SetCellValue(sheetName, cellName(row, col), club.ScoreForLossesInOT())

		col = colLetters[slices.Index(colLetters, col)+1]
		f.SetCellValue(sheetName, cellName(row, col), club.Score())

		row += 1
	}

	// Players
	for _, player := range team.Players {
		col = colOrigin

		f.SetCellValue(sheetName, cellName(row, col), player.Name)

		col = colLetters[slices.Index(colLetters, col)+1]
		f.SetCellValue(sheetName, cellName(row, col), player.ScoreForGoals())

		col = colLetters[slices.Index(colLetters, col)+1]
		f.SetCellValue(sheetName, cellName(row, col), player.ScoreForAssists())

		col = colLetters[slices.Index(colLetters, col)+2]
		f.SetCellValue(sheetName, cellName(row, col), player.Score())

		row += 1
	}

	style, err = f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		panic(err)
	}

	a := cellName(startRow, colLetters[slices.Index(colLetters, startCol)+1])
	b := cellName(row-1, col)
	f.SetCellStyle(sheetName, a, b, style)

	f.SetColWidth(sheetName, string(colLetters[slices.Index(colLetters, startCol)+1]), string(col), 7)
}

func classementSheet(f *excelize.File, sheetName string, curRanking, prevRanking Ranking) int {
	index, err := f.NewSheet(sheetName)
	if err != nil {
		panic(err)
	}

	// Header
	f.SetColWidth(sheetName, "A", "A", 5)

	f.SetColWidth(sheetName, "B", "B", 35)
	f.SetCellValue(sheetName, "B1", "Nom des equipes")

	f.SetCellValue(sheetName, "C1", "B/V")

	f.SetCellValue(sheetName, "D1", "Pass")

	f.SetCellValue(sheetName, "E1", "DP")

	f.SetCellValue(sheetName, "F1", "Points")

	f.SetCellValue(sheetName, "G1", "SEM.")

	// Populate the teams
	for i, t := range curRanking.Teams {
		deltaFromPrev := t.DeltaFrom(curRanking, prevRanking)

		produceRow(f, sheetName, i, t, deltaFromPrev)
	}

	numberOfTeams := len(curRanking.Teams)

	createCommentsBox(f, sheetName, numberOfTeams)

	centerColumnsCThroughG(f, sheetName, numberOfTeams)

	colorizeColumnF(f, sheetName, numberOfTeams)

	return index
}

func createCommentsBox(f *excelize.File, sheetName string, numberOfTeams int) {
	// Blank line
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", numberOfTeams+2), " ")
	f.MergeCell(sheetName, fmt.Sprintf("A%d", numberOfTeams+2), fmt.Sprintf("G%d", numberOfTeams+2))

	// Comments box
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", numberOfTeams+3), " ")
	f.SetRowHeight(sheetName, numberOfTeams+3, 100)
	f.MergeCell(sheetName, fmt.Sprintf("A%d", numberOfTeams+3), fmt.Sprintf("G%d", numberOfTeams+3))

	// Set border around cells
	style, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
		Alignment: &excelize.Alignment{
			Vertical: "top",
		},
	})
	if err != nil {
		panic(err)
	}

	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", numberOfTeams+3), fmt.Sprintf("G%d", numberOfTeams+3), style)
}

func centerColumnsCThroughG(f *excelize.File, sheetName string, numberOfTeams int) {
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		panic(err)
	}

	f.SetCellStyle(sheetName, "C1", fmt.Sprintf("G%d", numberOfTeams+2), style)
}

func colorizeColumnF(f *excelize.File, sheetName string, numberOfTeams int) {
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#0000FF",
			Bold:  true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		panic(err)
	}

	f.SetCellStyle(sheetName, "F2", fmt.Sprintf("F%d", numberOfTeams+1), style)
}

func produceRow(f *excelize.File, sheet string, n int, teamRank Rank, delta DeltaFromPrev) {
	n += 2

	// Position
	f.SetCellValue(sheet, "A"+strconv.Itoa(n), fmt.Sprintf("%02d-", n-1))

	// Team name
	f.SetCellValue(sheet, "B"+strconv.Itoa(n), fmt.Sprintf("%s (%s)", teamRank.Team.Name, teamRank.Team.Manager))

	// Goal points
	f.SetCellValue(sheet, "C"+strconv.Itoa(n), teamRank.Team.ScoreForGoals()+teamRank.Team.ScoreForWins())

	// Assist points
	f.SetCellValue(sheet, "D"+strconv.Itoa(n), teamRank.Team.ScoreForAssists())

	// Losses in OT points
	f.SetCellValue(sheet, "E"+strconv.Itoa(n), teamRank.Team.ScoreForLossesInOT())

	// Total score
	f.SetCellValue(sheet, "F"+strconv.Itoa(n), teamRank.Score)

	var deltaPos string
	if delta.Position > 0 {
		deltaPos = fmt.Sprintf("+%d", delta.Position)
	} else if delta.Position == 0 {
		deltaPos = "="
	} else {
		deltaPos = strconv.Itoa(delta.Position)
	}

	// Weekly delta
	f.SetCellValue(sheet, "G"+strconv.Itoa(n), fmt.Sprintf("%d/%s", delta.Score, deltaPos))
}
