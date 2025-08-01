package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
	"github.com/jdetok/golib/pgresd"
)

type Conf struct {
	l  logd.Logger
	db *sql.DB
	rc int64 // row counter
}

func main() {
	e := errd.InitErr()
	var sTime time.Time = time.Now()
	var cnf Conf
	// initialize logger
	l, err := logd.InitLogger("log", "etl")
	if err != nil {
		e.Msg = "error initializing logger"
		log.Fatal(e.BuildErr(err))
	}
	cnf.l = l
	var st string = "2004"
	var en string = "2007"
	szns, err := SznSlice(l, st, en)
	if err != nil {
		e.Msg = "error making seasons string"
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
	cnf.db = db

	cnf.rc = 0
	for _, s := range szns {
		sra := cnf.rc
		err = GLogSeasonETL(&cnf, s)
		if err != nil {
			e.Msg = "error inserting data"
			cnf.l.WriteLog(e.Msg)
			log.Fatal(e.BuildErr(err))
		}
		cnf.l.WriteLog(fmt.Sprint(
			"====  finished with ", s,
			fmt.Sprintf(
				"\n== total rows before: %d | total rows after: %d",
				sra, cnf.rc),
			fmt.Sprintf(
				"\n== rows affected from %s fetch: %d", s, cnf.rc-sra),
			fmt.Sprintf(
				"\n== total rows affected: %d", cnf.rc)))
	}

	cnf.l.WriteLog(fmt.Sprintf(
		"\n====  finished %d seasons between %s and %s | total rows affected: %d",
		len(szns), st, en, cnf.rc,
	))

	// send log file
	EmailLog(cnf.l)
	if err != nil {
		e.Msg = "error emailing log"
		cnf.l.WriteLog(e.Msg)
		log.Fatal(e.BuildErr(err))
	}

	var cTime time.Time = time.Now()
	cnf.l.WriteLog(
		fmt.Sprint(
			"process complete",
			fmt.Sprintf(
				"\n ---- start time: %v", sTime),
			fmt.Sprintf(
				"\n ---- cmplt time: %v", cTime),
			fmt.Sprintf(
				"\n ---- duration: %v", time.Since(sTime)),
			fmt.Sprintf(
				"\n---- etl for %d seasons between %s and %s | total rows affected: %d",
				len(szns), st, en, cnf.rc,
			),
		),
	)
}

/*
	if err := CrntPlayersETL(l, db, "1"); err != nil {
		e.Msg = "error getting players"
		l.WriteLog(e.Msg)
		log.Fatal(e.BuildErr(err))
	}
*/
// fetch & insert current (as of yesterday) stats for NBA and WNBA
// err = GLogDailyETL(l, db)

// err = GLogSeasonETL(cnf, SZN)
// if err != nil {
// 	e.Msg = "error inserting data"
// 	l.WriteLog(e.Msg)
// 	log.Fatal(e.BuildErr(err))
// }

// }
