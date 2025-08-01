package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
	"github.com/jdetok/golib/pgresd"
)

/* TODO -
* move the insert.go funcs to postgres package
* make chunk inserts concurrent
 */

// var SZN string = "2024-25"

// var seasons = []string{"2019-20", "2020-21", "2021-22", "2022-21", "2022-23", "2023-24"}

func SznSlice(l logd.Logger, start, end string) ([]string, error) {
	e := errd.InitErr()
	startYr, errS := strconv.Atoi(start)
	endYr, errE := strconv.Atoi(end)
	numY := endYr - startYr

	if errS != nil || errE != nil {
		e.Msg = "error converting start or end year to int"
		l.WriteLog(e.Msg)
		return nil, e.NewErr()
	}

	var szns []string
	for y := range numY {
		szns = append(szns,
			fmt.Sprintf(
				"%d-%s", startYr+y, strconv.Itoa(startYr + (y + 1))[2:]),
		)
	}
	szns = append(szns, fmt.Sprintf("%d-%s", endYr, strconv.Itoa(endYr + 1)[2:]))
	return szns, nil
}

/*
TODO AFTER MEDS:
create a type to store the logger, database, AND a row counter
update IN THE INSERT FUNC directly with pointer
*/

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
	var st string = "1994"
	var en string = "1996"
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
		err = GLogSeasonETL(&cnf, s)
		if err != nil {
			e.Msg = "error inserting data"
			cnf.l.WriteLog(e.Msg)
			log.Fatal(e.BuildErr(err))
		}
	}

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
				"\n ---- complete time: %v", cTime),
			fmt.Sprintf(
				"\n ---- duration: %v", time.Since(sTime)),
			fmt.Sprintf(
				"\n---- %d seasons between  %s and %s | total rows affected: %d",
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
