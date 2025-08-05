package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
)

// Conf struct, only have to pass this to access logger, db, row count, etc
type Conf struct {
	l    logd.Logger
	db   *sql.DB
	rc   int64 // row counter
	errs []string
}

func RunNightlyETL(cnf Conf) error {
	e := errd.InitErr()

	if err := CrntPlayersETL(cnf); err != nil {
		e.Msg = "error with current players ETL"
		cnf.l.WriteLog(e.Msg)
		return e.BuildErr(err)
	}

	if err := GLogDailyETL(&cnf); err != nil {
		e.Msg = "error with nightly game log ETL"
		cnf.l.WriteLog(e.Msg)
		return e.BuildErr(err)
	}

	cnf.l.WriteLog(fmt.Sprintf(
		"\n====  finished with nightly ETL | total rows affected: %d", cnf.rc))
	return nil
}

func RunSeasonETL(cnf Conf, startY, endY string) error {
	e := errd.InitErr()

	szns, err := SznBSlice(cnf.l, startY, endY)
	if err != nil {
		e.Msg = "error making seasons string"
		cnf.l.WriteLog(e.Msg)
		log.Fatal(e.BuildErr(err))
	}

	for _, s := range szns {
		sra := cnf.rc // capture row count at start of each season
		stT := time.Now()

		// players etl for season
		if err := SznPlayersETL(cnf, "1", s); err != nil {
			e.Msg = fmt.Sprint("error getting players for ", s)
			cnf.l.WriteLog(e.Msg)
			fmt.Println(e.BuildErr(err))
		}

		// get team and player game logs for the season
		err = GLogSeasonETL(&cnf, s)
		if err != nil {
			e.Msg = fmt.Sprint("error inserting data for ", s)
			cnf.l.WriteLog(e.Msg)
			fmt.Println(e.BuildErr(err))
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
		len(szns), startY, endY, cnf.rc,
	))

	return nil
}
