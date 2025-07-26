package main

import (
	"fmt"
	"log"
)

func main() {

	resp, err := RequestResp(nightlyPlGameLog)
	if err != nil {
		log.Fatalf("error getting response: %e", err)
	}
	fmt.Println(resp.ResultSets[0].Headers)
	for _, h := range resp.ResultSets[0].Headers {
		fmt.Println(h)
	}

	var intakePlGame = InsertStatement{
		Tbl:  "intake.gm_player",
		Cols: resp.ResultSets[0].Headers,
		Vals: resp.ResultSets[0].RowSet,
	}

	fmt.Println(len(intakePlGame.Cols))
	insStmnt := intakePlGame.Build()
	fmt.Println(insStmnt)
	pg := GetEnvPG()
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		fmt.Printf("Error connecting to postgres: %e\n", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pining postgres after successful pg.Conn(): %e\n", err)
	}
	fmt.Println("Successfully connected to & pinged postgres")

	r, err := db.Exec(insStmnt, intakePlGame.FlattenVals()...)
	if err != nil {
		log.Fatalf("Failed to insert values: %e\n", err)
	}
	ra, _ := r.RowsAffected()
	fmt.Printf("Rows Affected: %d\n", ra)

	// fmt.Println(intakePlGame.Build())
	// for _, h := range resp.ResultSets[0].Headers {
	// 	fmt.Println(h)
	// }

	// ProcessResp(resp)
}
