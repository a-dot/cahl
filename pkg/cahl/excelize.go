package cahl

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

const SHEET_NAME = "Classement"

func Excelize(curRanking, prevRanking Ranking, outputFile string) {
	f := excelize.NewFile()

	defer f.Close()

	index, err := f.NewSheet(SHEET_NAME)
	if err != nil {
		panic(err)
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Header
	f.SetColWidth(SHEET_NAME, "A", "A", 5)

	f.SetColWidth(SHEET_NAME, "B", "B", 35)
	f.SetCellValue(SHEET_NAME, "B1", "Nom des equipes")

	f.SetCellValue(SHEET_NAME, "C1", "B/V")

	f.SetCellValue(SHEET_NAME, "D1", "Pass")

	f.SetCellValue(SHEET_NAME, "E1", "DP")

	f.SetCellValue(SHEET_NAME, "F1", "Points")

	f.SetCellValue(SHEET_NAME, "G1", "SEM.")

	// Populate the teams
	for i, t := range curRanking.Teams {
		deltaFromPrev := t.DeltaFrom(curRanking, prevRanking)

		produceRow(f, SHEET_NAME, i, t, deltaFromPrev)
	}

	numberOfTeams := len(curRanking.Teams)

	createCommentsBox(f, numberOfTeams)

	centerColumnsCThroughG(f, numberOfTeams)

	colorizeColumnF(f, numberOfTeams)

	// Save spreadsheet by the given path.
	if err := f.SaveAs(outputFile); err != nil {
		panic(err)
	}
}

func createCommentsBox(f *excelize.File, numberOfTeams int) {
	// Blank line
	f.SetCellValue(SHEET_NAME, fmt.Sprintf("A%d", numberOfTeams+2), " ")
	f.MergeCell(SHEET_NAME, fmt.Sprintf("A%d", numberOfTeams+2), fmt.Sprintf("G%d", numberOfTeams+2))

	// Comments box
	f.SetCellValue(SHEET_NAME, fmt.Sprintf("A%d", numberOfTeams+3), " ")
	f.SetRowHeight(SHEET_NAME, numberOfTeams+3, 100)
	f.MergeCell(SHEET_NAME, fmt.Sprintf("A%d", numberOfTeams+3), fmt.Sprintf("G%d", numberOfTeams+3))

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

	f.SetCellStyle(SHEET_NAME, fmt.Sprintf("A%d", numberOfTeams+3), fmt.Sprintf("G%d", numberOfTeams+3), style)
}

func centerColumnsCThroughG(f *excelize.File, numberOfTeams int) {
	style, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		panic(err)
	}

	f.SetCellStyle(SHEET_NAME, "C1", fmt.Sprintf("G%d", numberOfTeams+2), style)
}

func colorizeColumnF(f *excelize.File, numberOfTeams int) {
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

	f.SetCellStyle(SHEET_NAME, "F2", fmt.Sprintf("F%d", numberOfTeams+1), style)
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
