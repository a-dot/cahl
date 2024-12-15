package main

const SHEET_NAME = "Classement"

// func main() {
// 	f := excelize.NewFile()

// 	defer f.Close()

// 	index, err := f.NewSheet(SHEET_NAME)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Set active sheet of the workbook.
// 	f.SetActiveSheet(index)
// 	f.DeleteSheet("Sheet1")

// 	// Header
// 	f.SetColWidth(SHEET_NAME, "A", "A", 5)

// 	f.SetColWidth(SHEET_NAME, "B", "B", 35)
// 	f.SetCellValue(SHEET_NAME, "B1", "Nom des equipes")

// 	// f.SetCellValue(SHEET_NAME, "C1", "B/V") // TODO MISSING

// 	// f.SetCellValue(SHEET_NAME, "D1", "Pass") // TODO MISSING

// 	// f.SetCellValue(SHEET_NAME, "E1", "DP") // TODO MISSING

// 	f.SetCellValue(SHEET_NAME, "C1", "Points")

// 	f.SetCellValue(SHEET_NAME, "D1", "SEM.")

// 	// Populate the teams
// 	inData, err := os.ReadFile("../cahl/output_20241204.json")
// 	if err != nil {
// 		panic(err)
// 	}

// 	var inFile teams.Output
// 	err = json.Unmarshal(inData, &inFile)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for i, t := range inFile.Ranking {
// 		produceRow(f, SHEET_NAME, i, &t)
// 	}

// 	// Blank line
// 	f.SetCellValue(SHEET_NAME, fmt.Sprintf("A%d", len(inFile.Ranking)+2), " ")
// 	f.MergeCell(SHEET_NAME, fmt.Sprintf("A%d", len(inFile.Ranking)+2), fmt.Sprintf("E%d", len(inFile.Ranking)+2))

// 	// Comments box
// 	f.SetCellValue(SHEET_NAME, fmt.Sprintf("A%d", len(inFile.Ranking)+3), " ")
// 	f.SetRowHeight(SHEET_NAME, len(inFile.Ranking)+3, 100)
// 	f.MergeCell(SHEET_NAME, fmt.Sprintf("A%d", len(inFile.Ranking)+3), fmt.Sprintf("E%d", len(inFile.Ranking)+3))

// 	style, err := f.NewStyle(&excelize.Style{
// 		Border: []excelize.Border{
// 			{Type: "left", Color: "000000", Style: 2},
// 			{Type: "top", Color: "000000", Style: 2},
// 			{Type: "bottom", Color: "000000", Style: 2},
// 			{Type: "right", Color: "000000", Style: 2},
// 		},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	f.SetCellStyle(SHEET_NAME, fmt.Sprintf("A%d", len(inFile.Ranking)+3), fmt.Sprintf("E%d", len(inFile.Ranking)+3), style)

// 	// Save spreadsheet by the given path.
// 	if err := f.SaveAs("ACHL.xlsx"); err != nil {
// 		panic(err)
// 	}
// }

// func produceRow(f *excelize.File, sheet string, n int, teamRank *teams.Ranking) {
// 	n += 2

// 	f.SetCellValue(sheet, "A"+strconv.Itoa(n), fmt.Sprintf("%02d-", n-1))
// 	f.SetCellValue(sheet, "B"+strconv.Itoa(n), teamRank.TeamName) // TODO Missing Manager's name from Team's name
// 	f.SetCellValue(sheet, "C"+strconv.Itoa(n), teamRank.Score)

// 	var deltaPos string
// 	if teamRank.Delta.Position > 0 {
// 		deltaPos = fmt.Sprintf("+%d", teamRank.Delta.Position)
// 	} else if teamRank.Delta.Position == 0 {
// 		deltaPos = "="
// 	} else {
// 		deltaPos = strconv.Itoa(teamRank.Delta.Position)
// 	}

// 	f.SetCellValue(sheet, "D"+strconv.Itoa(n), fmt.Sprintf("%d/%s", teamRank.Delta.Score, deltaPos))
// }
