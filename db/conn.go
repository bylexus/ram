package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var conn *sql.DB = nil

func Connect(path string) *sql.DB {
	var err error
	if conn == nil {
		conn, err = sql.Open("sqlite", path)
		if err != nil {
			panic(err)
		}
	}
	return conn

}

func GetConn() *sql.DB {
	if conn == nil {
		panic("DB not connected - use Connect() first")
	}
	return conn
}
