package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
)

// TODO: struct or something to standardize the lgs szn tables for a season etl

type LgTbls struct {
	lgs  []string
	tbls []Table
}

func PlayersParams() LgTbls {
	var lt LgTbls
	lt.lgs = []string{"00", "10"}
	lt.tbls = []Table{
		{
			Name:    "intake.gm_team",
			PrimKey: "game_id, team_id",
			PlTm:    "T",
		},
	}
	return lt
}

func GLogParams() LgTbls {
	var lt LgTbls
	lt.lgs = []string{"00", "10"}
	lt.tbls = []Table{
		{
			Name:    "intake.gm_team",
			PrimKey: "game_id, team_id",
			PlTm:    "T",
		},
		{
			Name:    "intake.gm_player",
			PrimKey: "game_id, player_id",
			PlTm:    "P",
		},
	}
	return lt
}

func GetPlayers(l logd.Logger, db *sql.DB, onlyCurrent string) error {
	e := errd.InitErr()
	sl := GetSeasons()
	var np = []string{"00", sl.Szn}
	var wp = []string{"10", sl.WSzn}
	var params = [][]string{np, wp}

	for _, p := range params {
		r := PlayerReq(onlyCurrent, p[0], p[1])
		resp, err := RequestResp(l, r)
		if err != nil {
			e.Msg = fmt.Sprintf("error getting response for %s", r.Endpoint)
			l.WriteLog(e.Msg)
			return e.BuildErr(err)
		}

		// get cols/rows from resp, return early when no rows in response
		var cols []string = resp.ResultSets[0].Headers
		var rows [][]any = resp.ResultSets[0].RowSet
		fmt.Println("Cols Length:", len(cols), "Rows Length:", len(rows))
		ProcessResp(resp)
	}

	return nil
}

// run single season
func GetManyGLogs(l logd.Logger, db *sql.DB, lgs []string, tbls []Table, szn string) error {
	e := errd.InitErr()
	for i := range lgs { // outer loop, 2 calls per lg
		for _, t := range tbls {
			// create request
			r := GameLogReq(lgs[i], szn, t.PlTm, "", "")
			l.WriteLog(fmt.Sprintf(
				"attempting to fetch %s: LG=%s, SZN=%s, PLTM=%s",
				r.Endpoint, lgs[i], szn, t.PlTm))

			// attempt to fetch & insert for current iteration
			err := GameLogETL(l, db, r, t.Name, t.PrimKey)
			if err != nil {
				e.Msg = fmt.Sprintf(
					"error during daily game log ETL. LG=%s, SZN=%s, PLTM=%s",
					lgs[i], szn, t.PlTm)
				l.WriteLog(e.Msg)
				return e.BuildErr(err)
			}
			// success, next call
			l.WriteLog(fmt.Sprintf(
				"finished with LG=%s, SZN=%s, PLTM=%s",
				lgs[i], szn, t.PlTm))
		}
	}
	return nil
}

func GLogSeasonETL(l logd.Logger, db *sql.DB, szn string) error {
	e := errd.InitErr()
	lt := GLogParams()
	err := GetManyGLogs(l, db, lt.lgs, lt.tbls, szn)
	if err != nil {
		e.Msg = fmt.Sprintf("error running ETL for %s", szn)
		l.WriteLog(e.Msg)
		return e.BuildErr(err)
	}
	return nil
}

/*
nightly game log fetch both PlayerTeam=P & T and NBA and WNBA
using yeseterday's date as DateFrom/DateTo
*/
func GLogDailyETL(l logd.Logger, db *sql.DB) error {
	e := errd.InitErr()
	yesterday := Yesterday(time.Now())
	lt := GLogParams()
	sl := GetSeasons()
	var szns = []string{sl.Szn, sl.WSzn}

	// makes 4 calls to leaguegamelog endpoint
	for i := range lt.lgs { // outer loop, 2 calls per lg
		for _, t := range lt.tbls {
			// create request
			r := GameLogReq(lt.lgs[i], szns[i], t.PlTm, yesterday, yesterday)
			l.WriteLog(fmt.Sprintf(
				"attempting to fetch %s: LG=%s, SZN=%s, PLTM=%s, DATE=%s",
				r.Endpoint, lt.lgs[i], szns[i], t.PlTm, yesterday))

			// attempt to fetch & insert for current iteration
			err := GameLogETL(l, db, r, t.Name, t.PrimKey)
			if err != nil {
				e.Msg = fmt.Sprintf(
					"error during daily game log ETL. LG=%s, SZN=%s, PLTM=%s, DATE=%s",
					lt.lgs[i], szns[i], t.PlTm, yesterday)
				l.WriteLog(e.Msg)
				return e.BuildErr(err)
			}
			// success, next call
			l.WriteLog(fmt.Sprintf(
				"finished with LG=%s, SZN=%s, PLTM=%s, DATE=%s",
				lt.lgs[i], szns[i], t.PlTm, yesterday))
		}
	}
	return nil
}

func GameLogETL(l logd.Logger, db *sql.DB, r GetReq, tbl, primKey string) error {
	e := errd.InitErr()

	// call endpoint in HTTP request, return Resp struct
	resp, err := RequestResp(l, r)
	if err != nil {
		e.Msg = fmt.Sprintf("error getting response for %s", r.Endpoint)
		l.WriteLog(e.Msg)
		return e.BuildErr(err)
	}

	// get cols/rows from resp, return early when no rows in response
	var cols []string = resp.ResultSets[0].Headers
	var rows [][]any = resp.ResultSets[0].RowSet
	if len(rows) == 0 {
		l.WriteLog("response returned 0 rows, exiting")
		return nil
	}
	l.WriteLog(
		fmt.Sprintf("response returned %d fields & %d rows",
			len(cols), len(rows)))

	// prepare the sql statement & chunks of values
	ins := MakeInsert(
		tbl,
		primKey,
		cols,
		rows,
	) // attempt to insert rows from response
	return ins.Insert(l, db)
}
