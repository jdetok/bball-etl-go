package main

import (
	"log"
	"time"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
	"github.com/jdetok/golib/maild"
	"github.com/jdetok/golib/pgresd"
)

var YESTERDAY string = time.Now().Add(-24 * time.Hour).Format("01/02/2006")

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

	if err := TeamSeasonRun(l, db, "00", "2016-17"); err != nil {
		e.Msg = "error running team season:"
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
