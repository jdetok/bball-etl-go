package main

import (
	"database/sql"
	"fmt"
	"log"
)

func BballETL(db *sql.DB, r GetReq, tbl string, primKey string) {
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
	fmt.Printf("attempting to insert data from %s into %s...\n", r.Endpoint, tbl)
	insStmnt := insert.Build()
	res, err := db.Exec(insStmnt, insert.FlattenVals()...)
	if err != nil {
		log.Fatalf("Failed to insert values: %e\n", err)
	}
	ra, _ := res.RowsAffected()
	fmt.Printf("%d Rows Affected: %s\n", ra, tbl)
}

/* EXAMPLES
BballETL(nightlyPlGameLog, "intake.gm_player", "game_id, player_id")
BballETL(nightlyTmGameLog, "intake.gm_team", "game_id, team_id")
BballETL(MakeGameLogReq("00", "2024-25", "T", "10/20/2024", "12/31/2024"),
		"intake.gm_team", "game_id, team_id")
BballETL(db, MakeGameLogReq("10", "2022-23", "T", "", ""),
		"intake.gm_team", "game_id, team_id")
*/
