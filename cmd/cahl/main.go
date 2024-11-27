package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/jessevdk/go-flags"

	"cahl/pkg/teams"
)

type Ranking struct {
	TeamName string        `json:"team"`
	Score    int           `json:"score"`
	Delta    *RankingDelta `json:"delta,omitempty"`
}

type RankingDelta struct {
	Score    int `json:"score"`
	Position int `json:"position"`
}

var opts struct {
	TeamsFile      string `short:"t" description:"teams file" default:"https://raw.githubusercontent.com/a-dot/cahl-teams/refs/heads/main/teams.json"`
	Season         string `short:"s" description:"season (format is YYYYXXXX)" default:"20242025"`
	DataOutputFile string `short:"d" description:"output json file with information used to calculate ranking"`
	Delta          string `short:"D" description:"calculate the delta from the last run by passing the output file here"`
}

type Output struct {
	Timestamp int64        `json:"timestamp"`
	Ranking   []Ranking    `json:"ranking"`
	Teams     []teams.Team `json:"teams"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	teams := teams.FromFile(opts.TeamsFile)

	ranking := make([]Ranking, 0, len(teams))

	for _, team := range teams {
		team.PopulatePlayersStats(opts.Season)

		team.PopulateClubsStats(opts.Season)

		score := team.Score()

		ranking = append(ranking, Ranking{
			TeamName: team.Name,
			Score:    score,
			Delta:    nil,
		})

		slog.Debug("total score for team", "team", team.Name, "score", score)
	}

	slices.SortFunc(ranking, func(a, b Ranking) int {
		return b.Score - a.Score
	})

	slog.Debug("ranking", "ranking", ranking)

	if len(opts.Delta) > 0 {
		populateDelta(opts.Delta, ranking)
	}

	if len(opts.DataOutputFile) > 0 {
		outputData, err := json.Marshal(Output{
			Timestamp: time.Now().Unix(),
			Ranking:   ranking,
			Teams:     teams,
		})
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(opts.DataOutputFile, outputData, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func findTeam(src []Ranking, team string) int {
	for i := range src {
		if src[i].TeamName == team {
			return i
		}
	}

	return -1
}

func populateDelta(prevFile string, current []Ranking) {
	prev := prevOutput(prevFile).Ranking

	for i, t := range current {
		prevTeamIdx := findTeam(prev, t.TeamName)

		d := &RankingDelta{
			Position: prevTeamIdx - i,
			Score:    t.Score - prev[prevTeamIdx].Score,
		}

		current[i].Delta = d
	}
}

func prevOutput(prevFile string) Output {
	data, err := os.ReadFile(prevFile)
	if err != nil {
		panic(err)
	}

	var r Output

	err = json.Unmarshal(data, &r)
	if err != nil {
		panic(err)
	}

	return r
}
