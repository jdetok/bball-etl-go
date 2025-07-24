package main

import (
	"fmt"
	"log"
)

var commonPlayerInfo = GetReq{
	Host:     HOST,
	Endpoint: "/stats/commonplayerinfo",
	Headers:  HDRS,
	Params:   []Pair{{"LeagueID", "10"}, {"PlayerID", "2544"}},
}

var playerAwards = GetReq{
	Host:     HOST,
	Endpoint: "/stats/playerawards",
	Headers:  HDRS,
	Params:   []Pair{{"PlayerID", "2544"}},
}

var commonAllPlayers = GetReq{
	Host:     HOST,
	Endpoint: "/stats/commonallplayers",
	Headers:  HDRS,
	Params: []Pair{
		{"LeagueID", "00"},
		{"IsOnlyCurrentSeason", "1"},
		{"Season", "2024-25"}},
}

func main() {
	// body, _, err := commonPlayerInfo.GetRespBody()
	// body, _, err := commonAllPlayers.GetRespBody()
	body, _, err := playerAwards.GetRespBody()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
