package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
)

func BballETL(l logd.Logger, db *sql.DB, r GetReq, tbl string, primKey string) {
	resp, err := RequestResp(r)
	if err != nil {
		log.Fatalf("error getting response: %e", err)
	}
	var insert = InsertStatement{
		Tbl:     tbl,
		PrimKey: primKey,
		Cols:    resp.ResultSets[0].Headers,
		Vals:    resp.ResultSets[0].RowSet,
	}

	l.WriteLog(fmt.Sprintf(
		"attempting to insert data from %s into %s...", r.Endpoint, tbl))

	insStmnt := insert.Build()
	res, err := db.Exec(insStmnt, insert.FlattenVals()...)
	if err != nil {
		log.Fatalf("Failed to insert values: %e\n", err)
	}
	ra, _ := res.RowsAffected()
	l.WriteLog(fmt.Sprintf(
		"insert statement into %s executed, Rows Affected: %d", tbl, ra))
}

func TeamSeasonRun(l logd.Logger, db *sql.DB, league, season string) error {
	e := errd.InitErr()

	var y1 string = season[0:4]
	y1int, err := strconv.Atoi(y1)
	if err != nil {
		e.Msg = "error converting year to int"
		return e.BuildErr(err)
	}
	var y2 string = strconv.Itoa(y1int + 1)

	var d1 = []string{fmt.Sprintf("10/20/%s", y1), fmt.Sprintf("12/31/%s", y1)}
	var d2 = []string{fmt.Sprintf("01/01/%s", y2), fmt.Sprintf("05/01/%s", y2)}
	var dates = [][]string{d1, d2}

	for _, d := range dates {
		BballETL(l, db, MakeGameLogReq(
			league, season, "T", d[0], d[1]),
			"intake.gm_team", "game_id, team_id")
	}
	return nil
}
