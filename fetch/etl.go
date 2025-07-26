package main

import (
	"fmt"
	"log"
)

func BballETL(r GetReq, tbl string, primKey string) {
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
	insStmnt := insert.Build()
	pg := GetEnvPG()
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		fmt.Printf("Error connecting to postgres: %e\n", err)
	}

	res, err := db.Exec(insStmnt, insert.FlattenVals()...)
	if err != nil {
		log.Fatalf("Failed to insert values: %e\n", err)
	}
	ra, _ := res.RowsAffected()
	fmt.Printf("%d Rows Affected: %s\n", ra, tbl)
}
