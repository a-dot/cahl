package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jessevdk/go-flags"

	"cahl/pkg/cahl"
	"cahl/pkg/nhlapi"
)

var opts struct {
	TeamsFile      string `short:"t" description:"teams file" default:"https://raw.githubusercontent.com/a-dot/cahl-teams/refs/heads/main/teams.json"`
	Season         string `short:"s" description:"season (format is YYYYXXXX)" default:"20242025"`
	DataOutputFile string `short:"d" description:"output json file with information used to calculate ranking"`
	Delta          string `short:"D" description:"calculate the delta from the last run by passing the output file here"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
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

	// if len(opts.Delta) > 0 {
	// 	populateDelta(opts.Delta, ranking)
	// }

	fmt.Println(cahl.CreateRanking(inTeams))

	// if len(opts.DataOutputFile) > 0 {
	// 	outputData, err := json.Marshal(teams.Output{
	// 		Timestamp: time.Now().Unix(),
	// 		Ranking:   ranking,
	// 		Teams:     inTeams,
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	err = os.WriteFile(opts.DataOutputFile, outputData, 0644)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}

// func findTeam(src []teams.Ranking, team string) int {
// 	for i := range src {
// 		if src[i].TeamName == team {
// 			return i
// 		}
// 	}

// 	return -1
// }

// func populateDelta(prevFile string, current []teams.Ranking) {
// 	prev := prevOutput(prevFile).Ranking

// 	for i, t := range current {
// 		prevTeamIdx := findTeam(prev, t.TeamName)

// 		d := &teams.RankingDelta{
// 			Position: prevTeamIdx - i,
// 			Score:    t.Score - prev[prevTeamIdx].Score,
// 		}

// 		current[i].Delta = d
// 	}
// }

// func prevOutput(prevFile string) teams.Output {
// 	data, err := os.ReadFile(prevFile)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var r teams.Output

// 	err = json.Unmarshal(data, &r)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return r
// }
