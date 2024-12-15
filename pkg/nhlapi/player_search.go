package nhlapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type PlayerSearch struct {
	cache map[string]uint64
}

func NewPlayerSearcher() PlayerSearch {
	return PlayerSearch{
		cache: make(map[string]uint64),
	}
}

type PlayerSearchResult struct {
	Id     string `json:"playerId"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
	Season string `json:"lastSeasonId"`
}

func (ps PlayerSearch) fetchPlayer(season, name string) (uint64, error) {
	ret, found := ps.cache[name+season]
	if found {
		return ret, nil
	}

	id, err := searchPlayerThroughAPI(season, name)
	if err != nil {
		return 0, err
	}

	ps.cache[name+season] = id

	return id, nil
}

func searchPlayerThroughAPI(season, name string) (uint64, error) {
	queryString := fmt.Sprintf("https://search.d3.nhle.com/api/v1/search/player?culture=en-us&limit=20&active=true&q=%s", url.QueryEscape(name))

	slog.Debug("searching player", "name", name, "queryString", queryString)

	resp, err := http.Get(queryString)
	if err != nil {
		return 0, fmt.Errorf("failed to contact nhl api, %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read nhl api server response, %w", err)
	}

	var psr []PlayerSearchResult

	err = json.Unmarshal(body, &psr)
	if err != nil {
		fmt.Println(string(body))
		panic(err)
	}

	for _, p := range psr {
		if p.Season != season || !p.Active {
			continue
		}

		if !checkName(name, p.Name) {
			continue
		}

		id, err := strconv.ParseUint(p.Id, 10, 64)
		if err != nil {
			slog.Error("error converting player ID to uint32", "id", p.Id)
			panic(err)
		}

		//TODO do we still need the suppliedID??
		// if suppliedID > 0 && suppliedID != id {
		// 	continue
		// }

		slog.Debug("search result for player", "name", name, "id", id)
		return id, nil
	}

	return 0, errors.New("search returned nothing")
}

func (ps PlayerSearch) Search(season, name string) (uint64, error) {
	//TODO do we still need the suppliedID??

	// var suppliedID uint64
	// rid := regexp.MustCompile(`(.*) +\(([0-9]*)\)`)
	// ridResult := rid.FindStringSubmatch(name)
	// if len(ridResult) > 1 {
	// 	name = ridResult[1]

	// 	suppliedID, err = strconv.ParseUint(ridResult[2], 10, 64)
	// 	if err != nil {
	// 		return 0, fmt.Errorf("error converting supplied player ID from name '%s', %w", name, err)
	// 	}
	// }

	return ps.fetchPlayer(season, name)
}

func checkName(name, searchResult string) bool {
	sname := strings.Split(name, " ")
	for _, n := range sname {
		if !strings.Contains(searchResult, n) {
			return false
		}
	}

	return true
}
