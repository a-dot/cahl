package cahl

import "fmt"

type PlayerSearcher interface {
	// Searches a player and returns his ID
	Search(season string, name string) (uint64, error)
}

type PlayerInfoFetcher interface {
	// Returns player's number of goals using his ID
	Goals(season string, id uint64) (int, error)

	// Returns player's number of assists using his ID
	Assists(season string, id uint64) (int, error)

	// Returns player's position (as string) using his ID
	Position(season string, id uint64) (string, error)
}

type ClubInfoFetcher interface {
	LossesOT(abbrev string) (int, error)
	Wins(abbrev string) (int, error)
}

func (p *Player) FetchStats(season string, s PlayerSearcher, info PlayerInfoFetcher) error {
	id, err := s.Search(season, p.Name)
	if err != nil {
		return fmt.Errorf("unable to find player, %w", err)
	}

	p.ID = id

	goals, err := info.Goals(season, p.ID)
	if err != nil {
		return fmt.Errorf("error fetching goals, %w", err)
	}
	p.Stats.Goals = goals

	assists, err := info.Assists(season, p.ID)
	if err != nil {
		return err
	}
	p.Stats.Assists = assists

	position, err := info.Position(season, p.ID)
	if err != nil {
		return fmt.Errorf("error fetching assists, %w", err)
	}

	p.Position, err = ParsePosition(position)
	if err != nil {
		return fmt.Errorf("error fetching position, %w", err)
	}

	return nil
}

func (c *Club) FetchStats(info ClubInfoFetcher) error {
	lossesOT, err := info.LossesOT(c.Abbrev)
	if err != nil {
		return fmt.Errorf("error fetching lossesOT, %w", err)
	}
	c.Stats.LossesInOT = lossesOT

	wins, err := info.Wins(c.Abbrev)
	if err != nil {
		return fmt.Errorf("error fetching wins, %w", err)
	}
	c.Stats.Wins = wins

	return nil
}
