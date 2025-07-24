package main

import (
	"fmt"
	"log"
)

var commonPlayerInfo = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/commonplayerinfo",
	Params:   []Pair{{"LeagueID", "10"}, {"PlayerID", "2544"}},
}

var playerAwards = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/playerawards",
	Params:   []Pair{{"PlayerID", "2544"}},
}

var commonAllPlayers = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/commonallplayers",
	Params: []Pair{
		{"LeagueID", "00"},
		{"IsOnlyCurrentSeason", "1"},
		{"Season", "2024-25"}},
}

var leagueStandings = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/leaguestandings",
	Params: []Pair{
		{"LeagueID", "00"},
		{"Season", "2024-25"},
		// {"SeasonType", "Playoffs"}},
		{"SeasonType", "Regular+Season"}},
}

func main() {
	// body, _, err := commonPlayerInfo.GetRespBody()
	// body, _, err := commonAllPlayers.GetRespBody()
	// body, _, err := playerAwards.GetRespBody()
	body, _, err := leagueStandings.GetRespBody()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
