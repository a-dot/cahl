package cahl

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

type TeamInput struct {
	Name    string   `json:"team_name"`
	Manager string   `json:"manager"`
	Players []string `json:"players"`
	Clubs   []string `json:"teams"`
}

func LoadTeams(inputFile string) ([]Team, error) {
	var data []byte
	var err error

	if strings.HasPrefix(inputFile, "http") {
		data, err = readFromRemote(inputFile)
	} else {
		data, err = readFromFile(inputFile)
	}

	if err != nil {
		return []Team{}, err
	}

	var it []TeamInput
	err = json.Unmarshal(data, &it)
	if err != nil {
		return []Team{}, err
	}

	return createTeams(it)
}

func readFromFile(fname string) ([]byte, error) {
	f, err := os.ReadFile(fname)
	if err != nil {
		return []byte{}, err
	}

	return f, nil
}

func readFromRemote(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func createTeams(in []TeamInput) ([]Team, error) {
	ret := make([]Team, 0, len(in))

	for _, inputTeam := range in {
		t := Team{
			Name:    inputTeam.Name,
			Manager: inputTeam.Manager,
			Players: make([]*Player, len(inputTeam.Players)),
			Clubs:   make([]*Club, len(inputTeam.Clubs)),
		}

		for i, p := range inputTeam.Players {
			t.Players[i] = &Player{
				Name: p,
			}
		}

		for i, c := range inputTeam.Clubs {
			t.Clubs[i] = &Club{
				Abbrev: c,
			}
		}

		ret = append(ret, t)
	}

	return ret, nil
}
