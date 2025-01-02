package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jessevdk/go-flags"

	"cahl/pkg/cahl"
	"cahl/pkg/nhlapi"
)

var opts struct {
	TeamsFile       string `short:"t" description:"teams file" default:"https://raw.githubusercontent.com/a-dot/cahl-teams/refs/heads/main/teams.json"`
	Season          string `short:"s" description:"season (format is YYYYXXXX)" default:"20242025"`
	DataOutputFile  string `short:"d" description:"output ranking in json format (that file is used to calculate ranking differential)"`
	PrevDataFile    string `short:"D" description:"calculate the delta from the last run by passing the output file here"`
	ExcelOutputFile string `short:"e" description:"excel output file name"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	var prevRanking cahl.Ranking
	if len(opts.PrevDataFile) > 0 {
		data, err := os.ReadFile(opts.PrevDataFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(data, &prevRanking)
		if err != nil {
			panic(err)
		}
	}

	inTeams, err := cahl.LoadTeams(opts.TeamsFile)
	if err != nil {
		panic(err)
	}

	playerSearcher := nhlapi.NewPlayerSearcher()
	playerFetcher := nhlapi.NewPlayerInfoFetcher()
	clubFetcher := nhlapi.NewClubInfoFetcher()

	for _, team := range inTeams {
		if err := team.Valid(); err != nil {
			slog.Error("invalid team", "name", team.Name, "err", err)
			os.Exit(1)
		}

		for _, player := range team.Players {
			err := player.FetchStats(opts.Season, playerSearcher, playerFetcher)
			if err != nil {
				panic(err)
			}
		}

		for _, club := range team.Clubs {
			err := club.FetchStats(clubFetcher)
			if err != nil {
				panic(err)
			}
		}
	}

	ranking := cahl.CreateRanking(inTeams)

	if len(opts.DataOutputFile) > 0 {
		outputData, err := json.Marshal(ranking)
		if err != nil {
			panic(err)
		}

		outputFile := fmt.Sprintf("%s_%s.json", opts.DataOutputFile, time.Now().Format("20060102"))

		err = os.WriteFile(outputFile, outputData, 0644)
		if err != nil {
			panic(err)
		}
	}

	if len(opts.ExcelOutputFile) > 0 {
		cahl.Excelize(inTeams, ranking, prevRanking, opts.ExcelOutputFile)
	}
}
