package main

import (
	"database/sql"
	"fmt"

	"github.com/jdetok/go-api-jdeko.me/getenv"
	_ "github.com/lib/pq"
)

type PostGres struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	ConnStr  string
}

func GetEnvPG() PostGres {
	var pg PostGres
	getenv.LoadDotEnv()
	pg.Host, _ = getenv.GetEnvStr("PG_HOST")
	pg.Port, _ = getenv.GetEnvInt("PG_PORT")
	pg.User, _ = getenv.GetEnvStr("PG_USER")
	pg.Password, _ = getenv.GetEnvStr("PG_PASS")
	pg.Database, _ = getenv.GetEnvStr("PG_DB")
	return pg
}

func (pg *PostGres) MakeConnStr() {
	pg.ConnStr = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pg.Host, pg.Port, pg.User, pg.Password, pg.Database)
}

func (pg *PostGres) Conn() (*sql.DB, error) {
	db, err := sql.Open("postgres", pg.ConnStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf(
			"Error pining postgres after successful conn: %e\n", err)
	}
	return db, err
}
