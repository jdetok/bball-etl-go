package main

import (
	"log"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
	"github.com/jdetok/golib/pgresd"
)

/* TODO -
* move the insert.go funcs to postgres package
* convert chunk insert to goroutine
 */

var SZN string = "2023-24"

func main() {
	e := errd.InitErr()

	// initialize logger
	l, err := logd.InitLogger("log", "etl")

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
	// err = GLogDailyETL(l, db)
	err = GLogSeasonETL(l, db, SZN)
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
