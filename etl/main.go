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
		e.Msg = "error connecting to postgres"
		l.WriteLog(e.Msg)
		log.Fatal(e.BuildErr(err))
	}

	// fetch & insert current (as of yesterday) stats for NBA and WNBA
	err = GLogDailyETL(l, db)
	if err != nil {
		e.Msg = "error inserting data"
		l.WriteLog(e.Msg)
		log.Fatal(e.BuildErr(err))
	}

	EmailLog(l.LogF)
	if err != nil {
		e.Msg = "error emailing log"
		l.WriteLog(e.Msg)
		log.Fatal(e.BuildErr(err))
	}
}