package main

func main() {
	// BballETL(nightlyPlGameLog, "intake.gm_player", "game_id, player_id")
	// BballETL(nightlyTmGameLog, "intake.gm_team", "game_id, team_id")
	// BballETL(MakeGameLogReq("00", "2024-25", "T", "10/20/2024", "12/31/2024"),
	// 	"intake.gm_team", "game_id, team_id")
	BballETL(MakeGameLogReq("00", "2024-25", "T", "01/01/2025", "06/30/2025"),
		"intake.gm_team", "game_id, team_id")
}
