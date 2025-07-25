package main

import (
	"encoding/json"
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
		{"Season", "2024-25"},
	},
}

var leagueStandings = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/leaguestandings",
	Params: []Pair{
		{"LeagueID", "00"},
		{"Season", "2024-25"},
		{"SeasonType", "Regular+Season"},
	},
}

var leagueGameLog = GetReq{
	Host:     HOST,
	Headers:  HDRS,
	Endpoint: "/stats/leaguegamelog",
	Params: []Pair{
		{"LeagueID", "00"},
		{"Season", "2024-25"},
		{"SeasonType", "Regular+Season"},
		{"Counter", "0"},
		{"PlayerOrTeam", "T"},
		{"Sorter", "DATE"},
		{"DateFrom", ""},
		{"DateTo", ""},
		{"Direction", "ASC"},
	},
}

func main() {
	// body, _, err := commonPlayerInfo.GetRespBody()
	// body, _, err := commonAllPlayers.GetRespBody()
	// body, _, err := playerAwards.GetRespBody()
	// body, _, err := leagueStandings.GetRespBody()
	body, _, err := leagueGameLog.GetRespBody()
	if err != nil {
		log.Fatalf("error getting response: %e", err)
	}
	var resp Resp
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Fatalf("error unmarshaling: %e", err)
	}
	fmt.Println(resp.ResultSets[0].RowSet[0]...)
	for _, r := range resp.ResultSets[0].RowSet {
		for i, x := range r {
			fmt.Printf("%v: %v\n", resp.ResultSets[0].Headers[i], x)
		}
	}
}
