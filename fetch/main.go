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
	body, _, err := commonAllPlayers.GetRespBody()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
