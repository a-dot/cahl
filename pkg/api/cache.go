package api

var playerCache map[string]uint64

var playerStatsCache map[uint64]*PlayerGoalsAssists

var teamCache map[string]TeamStats

func init() {
	playerCache = make(map[string]uint64)
	playerStatsCache = make(map[uint64]*PlayerGoalsAssists)
	teamCache = make(map[string]TeamStats)
}
