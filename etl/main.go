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

// Conf struct, only have to pass this to access logger, db, row count, etc
type Conf struct {
	l    logd.Logger
	db   *sql.DB
	rc   int64 // row counter
	errs []string
}

func main() {
	// start time variable for logging
	var sTime time.Time = time.Now()

	// SET START AND END SEASONS
	var st string = "1970"
	var en string = time.Now().Format("2006") // current year
	// var en string = "1970"

	// Conf variable, hold logger, db, etc
	var cnf Conf

	e := errd.InitErr() // start error handler

	// initialize logger
	l, err := logd.InitLogger("log", "full_etl")
	if err != nil {
		e.Msg = "error initializing logger"
		log.Fatal(e.BuildErr(err))
	}
	cnf.l = l // assign to cnf

	// postgres connection
	pg := pgresd.GetEnvPG()
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		e.Msg = "error connecting to postgres"
		cnf.l.WriteLog(e.Msg)
		log.Fatal(e.BuildErr(err))
	}

	cnf.db = db // asign to cnf
	cnf.db.SetMaxOpenConns(40)
	cnf.db.SetMaxIdleConns(20)
	// CREATE SLICE OF SEASONS FROM START/END YEARS
	szns, err := SznSlice(l, st, en)
	if err != nil {
		e.Msg = "error making seasons string"
		cnf.l.WriteLog(e.Msg)
		log.Fatal(e.BuildErr(err))
	}

	cnf.rc = 0 // START ROW COUNTER AT 0 BEFORE ETL STARTS
	// run ETL (http request, clean data, insert into db) for each season
	for _, s := range szns {
		sra := cnf.rc // capture row count at start of each season
		stT := time.Now()
		err = GLogSeasonETL(&cnf, s)
		if err != nil {
			e.Msg = "error inserting data"
			cnf.l.WriteLog(e.Msg)
			log.Fatal(e.BuildErr(err))
		} // log finished with season etl
		cnf.l.WriteLog(fmt.Sprint(
			fmt.Sprintf("====  finished with %s season ETL after %v",
				s, time.Since(stT)),
			fmt.Sprintf(
				"\n== total rows before: %d | total rows after: %d",
				sra, cnf.rc),
			fmt.Sprintf(
				"\n== rows affected from %s fetch: %d", s, cnf.rc-sra),
			fmt.Sprintf(
				"\n== total rows affected: %d", cnf.rc)))
	} // log finished with ETL
	cnf.l.WriteLog(fmt.Sprintf(
		"\n====  finished %d seasons between %s and %s | total rows affected: %d",
		len(szns), st, en, cnf.rc,
	))

	// write errors to the log
	if len(cnf.errs) > 0 {
		cnf.l.WriteLog(fmt.Sprintln("ERRORS:"))
		for _, e := range cnf.errs {
			cnf.l.WriteLog(fmt.Sprintln(e))
		}
	}

	// email log file to myself
	EmailLog(cnf.l)
	if err != nil {
		e.Msg = "error emailing log"
		cnf.l.WriteLog(e.Msg)
		log.Fatal(e.BuildErr(err))
	}

	// log process complete
	cnf.l.WriteLog(
		fmt.Sprint(
			"process complete",
			fmt.Sprintf(
				"\n ---- start time: %v", sTime),
			fmt.Sprintf(
				"\n ---- cmplt time: %v", time.Now()),
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
