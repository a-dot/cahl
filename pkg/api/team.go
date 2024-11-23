package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type defa struct {
	Value string `json:"default"`
}

type player struct {
	ID        int  `json:"id"`
	FirstName defa `json:"firstName"`
	LastName  defa `json:"lastName"`
}

type roster struct {
	Forwards   []player `json:"forwards"`
	DefenseMen []player `json:"defensemen"`
	Goalies    []player `json:"goalies"`
}

func PlayersFromTeam(teamAbbrev string) map[string]int {
	resp, err := http.Get(fmt.Sprintf("https://api-web.nhle.com/v1/roster/%s/current", teamAbbrev))
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var r roster
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err)
	}

	ret := make(map[string]int)

	for _, p := range r.Forwards {
		ret[fmt.Sprintf("%s%s", p.FirstName.Value, p.LastName.Value)] = p.ID
	}

	for _, p := range r.DefenseMen {
		ret[fmt.Sprintf("%s%s", p.FirstName.Value, p.LastName.Value)] = p.ID
	}

	for _, p := range r.Goalies {
		ret[fmt.Sprintf("%s%s", p.FirstName.Value, p.LastName.Value)] = p.ID
	}

	return ret
}

type TeamStats struct {
	Losses   int
	LossesOT int
	Wins     int
}

type Standing struct {
	TeamAbbrev struct {
		Abbrev string `json:"default"`
	} `json:"teamAbbrev"`

	Losses   int `json:"losses"`
	LossesOT int `json:"otLosses"`
	Ties     int `json:"ties"`
	Wins     int `json:"wins"`
}

type TeamStandings struct {
	Standings []Standing `json:"standings"`
}

func StatsForAllTeams() map[string]TeamStats {
	resp, err := http.Get("https://api-web.nhle.com/v1/standings/now")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var s TeamStandings

	err = json.Unmarshal(body, &s)
	if err != nil {
		panic(err)
	}

	ret := make(map[string]TeamStats)

	for _, v := range s.Standings {
		ret[v.TeamAbbrev.Abbrev] = TeamStats{
			Losses:   v.Losses,
			LossesOT: v.LossesOT,
			Wins:     v.Wins,
		}
	}

	return ret
}
