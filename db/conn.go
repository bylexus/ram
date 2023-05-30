package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var conn *sql.DB = nil

func Conn() *sql.DB {
	var err error
	if conn == nil {
		conn, err = sql.Open("sqlite", "test.db")
		if err != nil {
			panic(err)
		}
	}
	return conn
}
