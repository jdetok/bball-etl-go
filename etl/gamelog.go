package main

import (
	"fmt"
	"time"

	"github.com/jdetok/golib/errd"
)

func GameLogReqNew(league, season, sType, plTm, dateFrom, dateTo string) GetReq {
	var gr = GetReq{
		Host:     HOST,
		Headers:  HDRS,
		Endpoint: "/stats/leaguegamelog",
		Params: []Pair{
			{"LeagueID", league},
			{"Season", season},
			{"SeasonType", sType},
			{"Counter", "0"},
			{"Sorter", "DATE"},
			{"Direction", "DESC"},
			{"PlayerOrTeam", plTm},
			{"DateFrom", dateFrom},
			{"DateTo", dateTo},
		},
	}
	return gr
}

func GameLogReq(league, season, plTm, dateFrom, dateTo string) GetReq {
	var gr = GetReq{
		Host:     HOST,
		Headers:  HDRS,
		Endpoint: "/stats/leaguegamelog",
		Params: []Pair{
			{"LeagueID", league},
			{"Season", season},
			{"SeasonType", "Regular+Season"},
			{"Counter", "0"},
			{"Sorter", "DATE"},
			{"Direction", "DESC"},
			{"PlayerOrTeam", plTm},
			{"DateFrom", dateFrom},
			{"DateTo", dateTo},
		},
	}
	return gr
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

// run single season
func GetManyGLogs(cnf *Conf, lgs []string, tbls []Table, szn string) error {
	e := errd.InitErr()
	for i := range lgs { // outer loop, 2 calls per lg
		for _, t := range tbls {
			// get reg and playoffs
			for _, s := range []string{"Regular+Season", "Playoffs"} {
				// create request
				r := GameLogReqNew(lgs[i], szn, s, t.PlTm, "", "")
				cnf.l.WriteLog(fmt.Sprintf(
					"attempting to fetch %s: LG=%s, SZN=%s %s, PLTM=%s",
					r.Endpoint, lgs[i], szn, s, t.PlTm))

				// attempt to fetch & insert for current iteration
				err := GameLogETL(cnf, r, t.Name, t.PrimKey)
				if err != nil {
					e.Msg = fmt.Sprintf(
						"error during daily game log ETL. LG=%s, SZN=%s %s, PLTM=%s",
						lgs[i], szn, s, t.PlTm)
					cnf.l.WriteLog(e.Msg)
					return e.BuildErr(err)
				}
				// success, next call
				cnf.l.WriteLog(fmt.Sprintf(
					"finished with LG=%s, SZN=%s %s, PLTM=%s",
					lgs[i], szn, s, t.PlTm))
			}
		}
	}
	return nil
}

func GLogSeasonETL(cnf *Conf, szn string) error {
	e := errd.InitErr()
	lt := GLogParams()
	err := GetManyGLogs(cnf, lt.lgs, lt.tbls, szn)
	if err != nil {
		e.Msg = fmt.Sprintf("error running ETL for %s", szn)
		cnf.l.WriteLog(e.Msg)
		cnf.errs = append(cnf.errs, e.Msg) // capture if an error occured
		return e.BuildErr(err)
	}
	return nil
}

/*
nightly game log fetch both PlayerTeam=P & T and NBA and WNBA
using yeseterday's date as DateFrom/DateTo
*/
func GLogDailyETL(cnf *Conf) error {
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
			cnf.l.WriteLog(fmt.Sprintf(
				"attempting to fetch %s: LG=%s, SZN=%s, PLTM=%s, DATE=%s",
				r.Endpoint, lt.lgs[i], szns[i], t.PlTm, yesterday))

			// attempt to fetch & insert for current iteration
			err := GameLogETL(cnf, r, t.Name, t.PrimKey)
			if err != nil {
				e.Msg = fmt.Sprintf(
					"error during daily game log ETL. LG=%s, SZN=%s, PLTM=%s, DATE=%s",
					lt.lgs[i], szns[i], t.PlTm, yesterday)
				cnf.l.WriteLog(e.Msg)
				return e.BuildErr(err)
			}
			// success, next call
			cnf.l.WriteLog(fmt.Sprintf(
				"finished with LG=%s, SZN=%s, PLTM=%s, DATE=%s",
				lt.lgs[i], szns[i], t.PlTm, yesterday))
		}
	}
	return nil
}

func GameLogETL(cnf *Conf, r GetReq, tbl, primKey string) error {
	e := errd.InitErr()

	// call endpoint in HTTP request, return Resp struct
	resp, err := RequestResp(cnf.l, r)
	if err != nil {
		e.Msg = fmt.Sprintf("error getting response for %s", r.Endpoint)
		cnf.l.WriteLog(e.Msg)
		return e.BuildErr(err)
	}

	// get cols/rows from resp, return early when no rows in response
	var cols []string = resp.ResultSets[0].Headers
	var rows [][]any = resp.ResultSets[0].RowSet
	if len(rows) == 0 {
		cnf.l.WriteLog("response returned 0 rows, exiting")
		return nil
	}
	cnf.l.WriteLog(
		fmt.Sprintf("response returned %d fields & %d rows",
			len(cols), len(rows)))

	// prepare the sql statement & chunks of values
	ins := MakeInsert(
		tbl,
		primKey,
		cols,
		rows,
	) // attempt to insert rows from response
	return ins.InsertFast(cnf)
}
