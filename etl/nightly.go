package main

import (
	"fmt"

	"github.com/jdetok/golib/errd"
)

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
