package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type PlayerSearchResult struct {
	Id     string `json:"playerId"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
	Season string `json:"lastSeasonId"`
}

// SearchPlayer searches for a player's name and returns the player ID
// Optionally, the player's name can contain an ID that needs to match for the search to
// be successful. For example, "John Doe (1234)"
func SearchPlayer(name string, season string) uint64 {
	var err error

	var suppliedID uint64
	rid := regexp.MustCompile(`(.*) +\(([0-9]*)\)`)
	ridResult := rid.FindStringSubmatch(name)
	if len(ridResult) > 1 {
		name = ridResult[1]

		suppliedID, err = strconv.ParseUint(ridResult[2], 10, 64)
		if err != nil {
			slog.Error("error converting supplied player ID from name", "name", name)
			return 0
		}
	}

	queryString := fmt.Sprintf("https://search.d3.nhle.com/api/v1/search/player?culture=en-us&limit=20&q=%s", name)

	slog.Debug("searching player", "name", name, "queryString", queryString)

	resp, err := http.Get(queryString)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var psr []PlayerSearchResult

	err = json.Unmarshal(body, &psr)
	if err != nil {
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

		if suppliedID > 0 && suppliedID != id {
			continue
		}

		slog.Debug("search result for player", "name", name, "id", id)
		return id
	}

	return 0
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
