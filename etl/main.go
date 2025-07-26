package main

import (
	"fmt"
	"time"
)

var YESTERDAY string = time.Now().Add(-24 * time.Hour).Format("01/02/2006")

func main() {
	pg := GetEnvPG()
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		fmt.Printf("Error connecting to postgres: %e\n", err)
	}
	BballETL(db, MakeGameLogReq("10", "2025-26", "T", YESTERDAY, YESTERDAY),
		"intake.gm_team", "game_id, team_id")
	BballETL(db, MakeGameLogReq("10", "2025-26", "P", YESTERDAY, YESTERDAY),
		"intake.gm_player", "game_id, player_id")
}
