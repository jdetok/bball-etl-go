package main

import (
	"database/sql"
	"fmt"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
)

func BballETL(l logd.Logger, db *sql.DB, r GetReq, tbl, primKey string) error {
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

	ins := MakeInsert(
		tbl,
		primKey,
		cols,
		rows,
		// resp.ResultSets[0].Headers,
		// resp.ResultSets[0].RowSet,
	)
	return ins.Insert(l, db)
}
