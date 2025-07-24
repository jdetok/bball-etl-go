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

func main() {
	body, _, err := commonPlayerInfo.GetRespBody()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
