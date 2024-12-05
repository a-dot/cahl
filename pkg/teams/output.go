package teams

type Output struct {
	Timestamp int64     `json:"timestamp"`
	Ranking   []Ranking `json:"ranking"`
	Teams     []Team    `json:"teams"`
}
