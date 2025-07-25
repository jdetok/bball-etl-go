package main

import (
	"log"
)

func main() {
	resp, err := RequestResp(leagueGameLog)
	if err != nil {
		log.Fatalf("error getting response: %e", err)
	}
	ProcessResp(resp)
}
