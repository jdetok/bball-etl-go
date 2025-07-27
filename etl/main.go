package main

import (
	"log"
	"time"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
	"github.com/jdetok/golib/pgresd"
)

/* TODO -
* don't attempt to run the insert function if request returned no data

* convert the etl func to accept slice of str for league, season, pltm
** to enable calling both leagues, player and team fetch, multiple seasons, etc

* figure out how to request the appropriate season

* make a func to call BballETL for yesterday with the approriate league(s) and
** seasons

* probably move the insert.go funcs to postgres package
 */

var YESTERDAY string = Yesterday(time.Now())
var LEAGUE string = "00"
var NBA string = "00"
var WNBA string = "10"
var SEASON string = "2025-26"
var PLTM string = "T"
var DATEFROM string = ""
var DATETO string = ""

func main() {
	// e := errd.InitErr()
	// fmt.Println(CurrentSzns(time.Now()))

	// sl := GetSeasons()
	// fmt.Println(sl.Szn)
	// fmt.Println(sl.WSzn)

	// if err := RequestSchedule(SchedReq(WNBA, SEASON)); err != nil {
	// 	e.Msg = "req error"
	// 	log.Fatal(e.BuildErr(err))
	// }

	e := errd.InitErr()

	// initialize logger
	l, err := logd.InitLogger("log", "test")

	if err != nil {
		e.Msg = "error initializing logger"
		log.Fatal(e.BuildErr(err))
	}

	// postgres connection
	pg := pgresd.GetEnvPG()
	pg.MakeConnStr()
	db, err := pg.Conn()

	if err != nil {
		e.Msg = "error connecting to postgres:"
		log.Fatal(e.BuildErr(err))
	}

	// err = GameLogETL(l, db, GameLogReq(LEAGUE, SEASON, PLTM, DATEFROM, DATETO),
	// 	"intake.gm_team", "game_id, team_id")
	err = GLogDailyETL(l, db)
	if err != nil {
		e.Msg = "error inserting data"
		log.Fatal(e.BuildErr(err))
	}

	EmailLog(l.LogF)

	if err != nil {
		e.Msg = "error emailing log:"
		log.Fatal(e.BuildErr(err))
	}
}

// func EmailLog(file string) error {
// 	m := maild.MakeMail(
// 		[]string{"jdekock17@gmail.com"},
// 		"Go bball ETL log attached",
// 		"the Go bball ETL process ran. The log is attached.",
// 	)
// 	return m.SendMIMEEmail(file)
// }

/*
	// wnba team then player
	BballETL(l, db, MakeGameLogReq("10", "2025-26", "T", YESTERDAY, YESTERDAY),
		"intake.gm_team", "game_id, team_id")
	BballETL(l, db, MakeGameLogReq("10", "2025-26", "P", YESTERDAY, YESTERDAY),
		"intake.gm_player", "game_id, player_id")

*/
