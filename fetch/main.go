package main

import "fmt"

func main() {
	pg := GetEnvPG()
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		fmt.Printf("Error connecting to postgres: %e\n", err)
	}
	// BballETL(nightlyPlGameLog, "intake.gm_player", "game_id, player_id")
	// BballETL(nightlyTmGameLog, "intake.gm_team", "game_id, team_id")
	// BballETL(MakeGameLogReq("00", "2024-25", "T", "10/20/2024", "12/31/2024"),
	// 	"intake.gm_team", "game_id, team_id")
	BballETL(db, MakeGameLogReq("10", "2025-26", "T", "", ""),
		"intake.gm_team", "game_id, team_id")
}
