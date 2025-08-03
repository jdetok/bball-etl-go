package main

import (
	"fmt"

	"github.com/jdetok/golib/errd"
)

func PlayerReq(onlyCurrent, league, season string) GetReq {
	var gr = GetReq{
		Host:     HOST,
		Headers:  HDRS,
		Endpoint: "/stats/commonallplayers",
		Params: []Pair{
			{"IsOnlyCurrentSeason", onlyCurrent},
			{"LeagueID", league},
			{"Season", season},
		},
	}
	return gr
}

func PlayersParams() LgTbls {
	var lt LgTbls
	lt.lgs = []string{"00", "10"}
	lt.tbls = []Table{
		{
			Name:    "intake.player",
			PrimKey: "person_id",
		},
		{
			Name:    "intake.wplayer",
			PrimKey: "person_id",
		},
	}
	return lt
}

// SAME AS CURRENT PLAYER ETL BUT FOR INDIVIDUAL SEASON
// WILL NEED A NEW GET SEASONS FUNCTION AS WELL
func SznPlayersETL(cnf Conf, onlyCurrent, season string) error {
	e := errd.InitErr()
	pp := PlayersParams()

	cnf.l.WriteLog(fmt.Sprintf(
		"attempting players ETL for %s nba/wnba seasons",
		season))
	for i := range pp.lgs {
		var lg string
		switch pp.lgs[i] {
		case "00":
			lg = "nba"
		case "10":
			lg = "wnba"
		}

		cnf.l.WriteLog(fmt.Sprintf("attempting to insert %s %s players", season, lg))
		// r := PlayerReq(onlyCurrent, p[0], p[1])
		r := PlayerReq(onlyCurrent, pp.lgs[i], season)
		resp, err := RequestResp(cnf.l, r)
		if err != nil {
			e.Msg = fmt.Sprintf("error getting response for %s: lg: %s szn: %s", r.Endpoint, lg, season)
			cnf.l.WriteLog(e.Msg)
			return e.BuildErr(err)
		}

		// get cols/rows from resp, return early when no rows in response
		var cols []string = resp.ResultSets[0].Headers
		var rows [][]any = resp.ResultSets[0].RowSet
		// ProcessResp(resp)
		fmt.Println("Cols Length:", len(cols), "Rows Length:", len(rows))

		if len(rows) == 0 {
			cnf.l.WriteLog("response returned 0 rows, exiting")
			return nil
		}
		cnf.l.WriteLog(
			fmt.Sprintf("response returned %d fields & %d rows",
				len(cols), len(rows)))

		// prepare the sql statement & chunks of values
		ins := MakeInsert(
			pp.tbls[i].Name,
			pp.tbls[i].PrimKey,
			cols,
			rows,
		) // attempt to insert rows from response
		ins.Insert(&cnf)

		cnf.l.WriteLog(fmt.Sprintf("%s %s players ETL complete", season, lg))
	}
	cnf.l.WriteLog(fmt.Sprint("players ETL complete for ", season))
	return nil
}

func CrntPlayersETL(cnf Conf, onlyCurrent string) error {
	e := errd.InitErr()
	sl := GetSeasons()
	var szns = []string{sl.Szn, sl.WSzn}
	pp := PlayersParams()

	cnf.l.WriteLog(fmt.Sprintf(
		"attempting current players ETL for %s nba season and %s wnba season",
		sl.Szn, sl.WSzn))
	for i := range pp.lgs {
		var lg string
		switch pp.lgs[i] {
		case "00":
			lg = "nba"
		case "10":
			lg = "wnba"
		}

		cnf.l.WriteLog(fmt.Sprintf("attempting to insert current %s players", lg))
		// r := PlayerReq(onlyCurrent, p[0], p[1])
		r := PlayerReq(onlyCurrent, pp.lgs[i], szns[i])
		resp, err := RequestResp(cnf.l, r)
		if err != nil {
			e.Msg = fmt.Sprintf("error getting response for %s", r.Endpoint)
			cnf.l.WriteLog(e.Msg)
			return e.BuildErr(err)
		}

		// get cols/rows from resp, return early when no rows in response
		var cols []string = resp.ResultSets[0].Headers
		var rows [][]any = resp.ResultSets[0].RowSet
		// ProcessResp(resp)
		fmt.Println("Cols Length:", len(cols), "Rows Length:", len(rows))

		if len(rows) == 0 {
			cnf.l.WriteLog("response returned 0 rows, exiting")
			return nil
		}
		cnf.l.WriteLog(
			fmt.Sprintf("response returned %d fields & %d rows",
				len(cols), len(rows)))

		// prepare the sql statement & chunks of values
		ins := MakeInsert(
			pp.tbls[i].Name,
			pp.tbls[i].PrimKey,
			cols,
			rows,
		) // attempt to insert rows from response
		ins.Insert(&cnf)

		cnf.l.WriteLog(fmt.Sprintf("current %s players ETL complete", lg))
	}
	cnf.l.WriteLog("current players ETL complete for all leagues")
	return nil
}
