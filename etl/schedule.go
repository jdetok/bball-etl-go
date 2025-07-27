package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jdetok/golib/errd"
)

/*
call stats/scheduleleaguev2 to get current season
*/

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

type SznLg struct {
	League    string
	Season    string
	StartDate string
	EndDate   string
}

type SeasonLeague struct {
	Szn  string
	WSzn string
}

/*
returns slice of season strings for date (generally pass time.Now())
calling in 2025 will return 2024-25 and 2025-26 and so on
*/
func CurrentSzns(dt time.Time) []string {
	var cyyy string = dt.Format("2006")
	var cy string = dt.AddDate(1, 0, 0).Format("06")

	var pyyy string = dt.AddDate(-1, 0, 0).Format("2006")
	var py string = dt.Format("06")

	return []string{
		fmt.Sprint(pyyy, "-", py),
		fmt.Sprint(cyyy, "-", cy),
	}
}

// pass a time (usually time.Now()), return string with yesterday's date
func Yesterday(dt time.Time) string {
	return dt.Add(-24 * time.Hour).Format("01/02/2006")
}

/*
	func MakeSznLg(startDate, endDate string) SznLg {
		lgs := []string{"00", "10"}
		szns := []string{
			fmt.Sprintf("%s-%s",
				time.Now().Format("2006"),
			),
		}
		var sl = SznLg{
			League:    lg,
			Season:    szn,
			StartDate: startDate,
			EndDate:   endDate,
		}
	}
*/
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
