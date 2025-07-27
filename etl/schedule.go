package main

import (
	"encoding/json"
	"fmt"

	"github.com/jdetok/golib/errd"
)

type RespSched struct {
	Dates GameDates `json:"leagueSchedule"`
}

// main json object in response body after endpoint/params
type GameDates struct {
	GmDates []GameDate `json:"gameDates"`
}

type GameDate struct {
	Date string `json:"gameDate"`
}

func SchedReq(league, season string) GetReq {
	var gr = GetReq{
		Host:     HOST,
		Headers:  HDRS,
		Endpoint: "/stats/scheduleleaguev2",
		Params: []Pair{
			{"LeagueID", league},
			{"Season", season},
		},
	}
	return gr
}

func RequestSchedule(gr GetReq) error {
	e := errd.InitErr()
	fmt.Printf("requesting data from %s...\n", gr.Endpoint)
	body, err := gr.BodyFromReq()
	if err != nil {
		e.Msg = "error getting schedule response"
		return e.BuildErr(err)
	}

	var resp RespSched
	if err := json.Unmarshal(body, &resp); err != nil {
		e.Msg = "error unmarshaling schedule response"
		fmt.Println(err)
		return e.BuildErr(err)
	}

	fmt.Println(resp.Dates)
	return nil
}
