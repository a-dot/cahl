package nhlapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ClubInfo struct {
	cache map[string]*ClubStats
}

func NewClubInfoFetcher() ClubInfo {
	return ClubInfo{
		cache: make(map[string]*ClubStats),
	}
}

type ClubStats struct {
	Losses   int
	LossesOT int
	Wins     int
}

func (ci ClubInfo) LossesOT(abbrev string) (int, error) {
	cs, err := ci.fetchClubStats(abbrev)
	if err != nil {
		return 0, nil
	}

	return cs.LossesOT, nil
}

func (ci ClubInfo) Wins(abbrev string) (int, error) {
	cs, err := ci.fetchClubStats(abbrev)
	if err != nil {
		return 0, nil
	}

	return cs.Wins, nil
}

func (ci ClubInfo) fetchClubStats(abbrev string) (*ClubStats, error) {
	ret, found := ci.cache[abbrev]
	if found {
		return ret, nil
	}

	if err := ci.buildClubsCache(); err != nil {
		return nil, err
	}

	ret, found = ci.cache[abbrev]
	if !found {
		return nil, fmt.Errorf("invalid abbrev for club '%s'", abbrev)
	}

	return ret, nil
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

func (ci ClubInfo) buildClubsCache() error {
	resp, err := http.Get("https://api-web.nhle.com/v1/standings/now")
	if err != nil {
		return fmt.Errorf("failed to contact nhl api, %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read nhl api server response, %w", err)
	}

	var s TeamStandings

	err = json.Unmarshal(body, &s)
	if err != nil {
		panic(err)
	}

	for _, v := range s.Standings {
		ci.cache[v.TeamAbbrev.Abbrev] = &ClubStats{
			Losses:   v.Losses,
			LossesOT: v.LossesOT,
			Wins:     v.Wins,
		}
	}

	return nil
}
