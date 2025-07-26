package main

import (
	"database/sql"
	"fmt"
	"log"

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

	// fmt.Printf("%d Rows Affected: %s\n", ra, tbl)
	l.WriteLog(fmt.Sprintf(
		"insert statement into %s executed, Rows Affected: %d", tbl, ra))
}
