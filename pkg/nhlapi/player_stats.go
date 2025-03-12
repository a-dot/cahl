package nhlapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type PlayerInfo struct {
	cache map[uint64]*PlayerStats
}

func NewPlayerInfoFetcher() PlayerInfo {
	return PlayerInfo{
		cache: make(map[uint64]*PlayerStats),
	}
}

type PlayerStats struct {
	Position string
	Assists  int
	Goals    int
}

func (pi PlayerInfo) Goals(season string, id uint64) (int, error) {
	ps, err := pi.fetchPlayerStats(season, id)
	if err != nil {
		return 0, err
	}

	return ps.Goals, nil
}

func (pi PlayerInfo) Assists(season string, id uint64) (int, error) {
	ps, err := pi.fetchPlayerStats(season, id)
	if err != nil {
		return 0, err
	}

	return ps.Assists, nil
}

func (pi PlayerInfo) Position(season string, id uint64) (string, error) {
	ps, err := pi.fetchPlayerStats(season, id)
	if err != nil {
		return "", err
	}

	return ps.Position, nil
}

func (pi PlayerInfo) fetchPlayerStats(season string, id uint64) (*PlayerStats, error) {
	ret, found := pi.cache[id]
	if found {
		return ret, nil
	}

	stats, err := searchPlayerStatsThroughAPI(season, id)
	if err != nil {
		return nil, err
	}

	pi.cache[id] = stats

	return stats, nil
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

func searchPlayerStatsThroughAPI(season string, id uint64) (*PlayerStats, error) {
	seasonID, err := strconv.Atoi(season)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(fmt.Sprintf("https://api-web.nhle.com/v1/player/%d/landing", id))
	if err != nil {
		return nil, fmt.Errorf("failed to contact nhl api, %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read nhl api server response, %w", err)
	}

	var p PlayerStatsAllSeasons

	err = json.Unmarshal(body, &p)
	if err != nil {
		panic(err)
	}

	ret := &PlayerStats{}
	found := false

	for _, v := range p.Totals {
		if v.Season == seasonID {
			if ret.Position != "" && ret.Position != p.Position {
				return nil, fmt.Errorf("player (%d) had position '%s' and now we found position '%s'", id, ret.Position, p.Position)
			}

			ret.Position = p.Position
			ret.Goals += v.Goals
			ret.Assists += v.Assists

			found = true
		}
	}

	if !found {
		return nil, fmt.Errorf("player stats not found for player id %d", id)
	}

	return ret, nil
}
