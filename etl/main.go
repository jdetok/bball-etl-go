package main

import (
	"log"
	"time"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
	"github.com/jdetok/golib/maild"
	"github.com/jdetok/golib/pgresd"
)

/* TODO -
* don't attempt to run the insert function if request returned no data

* convert the make request func to accept slice of str for league, season, pltm
** to enable calling both leagues, player and team fetch, multiple seasons, etc

* figure out how to request the appropriate season
 */

var YESTERDAY string = time.Now().Add(-24 * time.Hour).Format("01/02/2006")
var LEAGUE string = "00"
var NBA string = "00"
var WNBA string = "10"
var SEASON string = "2011-12"
var PLTM string = "T"
var DATEFROM string = ""
var DATETO string = ""

func main() {
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

	err = BballETL(l, db, MakeGameLogReq(LEAGUE, SEASON, PLTM, DATEFROM, DATETO),
		"intake.gm_team", "game_id, team_id")
	if err != nil {
		e.Msg = "error inserting data"
		log.Fatal(e.BuildErr(err))
	}

	// send email with log attached
	EmailLog(l.LogF)
	if err != nil {
		e.Msg = "error emailing log:"
		log.Fatal(e.BuildErr(err))
	}
}

func EmailLog(file string) error {
	m := maild.MakeMail(
		[]string{"jdekock17@gmail.com"},
		"Go bball ETL log attached",
		"the Go bball ETL process ran. The log is attached.",
	)
	return m.SendMIMEEmail(file)
}

/*
	// wnba team then player
	BballETL(l, db, MakeGameLogReq("10", "2025-26", "T", YESTERDAY, YESTERDAY),
		"intake.gm_team", "game_id, team_id")
	BballETL(l, db, MakeGameLogReq("10", "2025-26", "P", YESTERDAY, YESTERDAY),
		"intake.gm_player", "game_id, player_id")

*/
