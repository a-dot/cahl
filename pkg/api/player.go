package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type PlayerGoalsAssists struct {
	Position string
	Assists  int
	Goals    int
}

type SeasonTotals struct {
	Assists int `json:"assists"`
	Goals   int `json:"goals"`
	Season  int `json:"season"`
}

type PlayerStatsAllSeasons struct {
	Position string         `json:"position"`
	Totals   []SeasonTotals `json:"seasonTotals"`
}

func PlayerStats(id uint64, season string) *PlayerGoalsAssists {
	seasonID, err := strconv.Atoi(season)
	if err != nil {
		panic(err)
	}

	ret, found := playerStatsCache[id]
	if found {
		return ret
	}

	resp, err := http.Get(fmt.Sprintf("https://api-web.nhle.com/v1/player/%d/landing", id))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var p PlayerStatsAllSeasons

	err = json.Unmarshal(body, &p)
	if err != nil {
		panic(err)
	}

	for _, v := range p.Totals {
		if v.Season == seasonID {
			ret = &PlayerGoalsAssists{
				Position: p.Position,
				Goals:    v.Goals,
				Assists:  v.Assists,
			}

			playerStatsCache[id] = ret

			return ret
		}
	}

	return nil
}
