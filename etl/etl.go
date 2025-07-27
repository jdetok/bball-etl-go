package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
)

/*
request
*/

type Table struct {
	Name    string
	PrimKey string
	PlTm    string
}

func GLogDailyETL(l logd.Logger, db *sql.DB) error {
	e := errd.InitErr()
	sl := GetSeasons()
	yesterday := Yesterday(time.Now())
	var lgs = []string{"00", "10"}
	var szns = []string{sl.Szn, sl.WSzn}
	var tbls = []Table{
		Table{
			Name:    "intake.gm_team",
			PrimKey: "game_id, team_id",
			PlTm:    "T",
		},
		Table{
			Name:    "intake.gm_player",
			PrimKey: "game_id, player_id",
			PlTm:    "P",
		},
	}

	for i := range lgs {
		for _, t := range tbls {
			err := GameLogETL(l, db, GameLogReq(
				lgs[i], szns[i], t.PlTm, yesterday, yesterday,
			), t.Name, t.PrimKey)
			if err != nil {
				e.Msg = fmt.Sprintf(
					"error during daily game log ETL. LG=%s, SZN=%s, PLTM=%s, DATE=%s",
					lgs[i], szns[i], t.PlTm, yesterday)
				return e.BuildErr(err)
			}
		}
	}
	return nil
}

func GameLogETL(l logd.Logger, db *sql.DB, r GetReq, tbl, primKey string) error {
	e := errd.InitErr()

	l.WriteLog(fmt.Sprintf("attempting to get data from %s", r.Endpoint))
	resp, err := RequestResp(r)
	if err != nil {
		e.Msg = fmt.Sprintf("error getting response for %s", r.Endpoint)
		return e.BuildErr(err)
	}

	var cols []string = resp.ResultSets[0].Headers
	var rows [][]any = resp.ResultSets[0].RowSet

	l.WriteLog(
		fmt.Sprintf("response returned %d fields & %d rows",
			len(cols), len(rows)))

	// return early when no rows in response
	if len(rows) == 0 {
		return nil
	}

	ins := MakeInsert(
		tbl,
		primKey,
		cols,
		rows,
	)
	return ins.Insert(l, db)
}
