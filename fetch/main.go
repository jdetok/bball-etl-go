package main

func main() {
	BballETL(nightlyPlGameLog, "intake.gm_player", "game_id, player_id")
	BballETL(nightlyTmGameLog, "intake.gm_team", "game_id, team_id")
}
