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

	var intakeTmGame = InsertStatement{
		Tbl:  "intake.gm_player",
		Cols: resp.ResultSets[0].Headers,
		Vals: resp.ResultSets[0].RowSet,
	}

	insStmnt := intakeTmGame.Build()
	fmt.Println(insStmnt)
	// db.Exec(insStmnt, intakeTmGame.FlattenVals())

	// fmt.Println(intakePlGame.Build())
	// for _, h := range resp.ResultSets[0].Headers {
	// 	fmt.Println(h)
	// }

	// ProcessResp(resp)
}
