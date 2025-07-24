package main

import (
	"fmt"
	"log"
)

var hdrs = []Pair{
	{"Accept", "application/json"},
	{"Connection", "keep-alive"},
	{"Referer", "https://www.nba.com"},
	{"Origin", "https://www.nba.com"},
	{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"},
}

var params = []Pair{
	{"LeagueID", "10"},
	{"PlayerID", "2544"},
}

var getReq = GetReq{
	Host:     "stats.nba.com",
	Endpoint: "/stats/commonplayerinfo",
	Params:   params,
	Headers:  hdrs,
}

func main() {
	body, _, err := getReq.GetRespBody()
	// body, _, err := Get("stats.nba.com", "/stats/commonplayerinfo", params, hdrs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
