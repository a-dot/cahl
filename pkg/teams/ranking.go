package teams

type Ranking struct {
	TeamName string        `json:"team"`
	Score    int           `json:"score"`
	Delta    *RankingDelta `json:"delta,omitempty"`
}

type RankingDelta struct {
	Score    int `json:"score"`
	Position int `json:"position"`
}
