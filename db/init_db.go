package db

import (
	"database/sql"
	"fmt"
	"sort"

	e "github.com/bylexus/go-stdlib/err"
	l "github.com/bylexus/go-stdlib/log"
	"github.com/bylexus/go-stdlib/maps"
	_ "modernc.org/sqlite"
)

// Add DB migrations here, in the correct order. They will be executed on the table
// if not yet applied.
var migrations map[int64]func(*sql.DB) error = map[int64]func(*sql.DB) error{
	1: dbMigration_00001,
	2: dbMigration_00002,
	3: dbMigration_00003,
}

func InitDb(logger *l.SeverityLogger, conn *sql.DB) {
	var err error

	logger.Info("Start DB Migrations")

	// PRAGMA user_version
	// see https://www.sqlite.org/pragma.html#pragma_user_version:
	// PRAGMA user_version can be used to set a version number in the db.
	// We use it to get the version number the schema is in, to run the apropriate
	// migrations.
	schemaVersion := getSchemaVersion(conn)
	logger.Info("DB Schema is in version %d", schemaVersion)

	availableVersions := maps.GetMapKeys(&migrations)
	sort.Slice(availableVersions, func(i int, j int) bool {
		return availableVersions[i] < availableVersions[j]
	})

	for _, version := range availableVersions {
		if version > schemaVersion {
			logger.Info("Executing DB Migration #%d ... ", version)
			err = migrations[version](conn)
			if err != nil {
				panic(err)
			}
			setSchemaVersion(conn, version)
			logger.Info("done")
		}
	}
}

func getSchemaVersion(conn *sql.DB) int64 {
	var version int64 = 0
	res := conn.QueryRow("PRAGMA user_version")
	res.Scan(&version)
	return version
}

func setSchemaVersion(conn *sql.DB, version int64) {
	_, err := conn.Exec(fmt.Sprintf("PRAGMA user_version = %d", version))
	e.PanicOnErr(err)
}

func dbMigration_00001(conn *sql.DB) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS note (
			id VARCHAR(36) NOT NULL PRIMARY KEY,
			note TEXT,
			url TEXT,
			tags TEXT,
			done BOOLEAN
		)
	`)
	return err
}

func dbMigration_00002(conn *sql.DB) error {
	_, err := conn.Exec(`
	ALTER TABLE note ADD COLUMN user_id BIGINT
	`)
	return err
}

func dbMigration_00003(conn *sql.DB) error {
	_, err := conn.Exec(`
	ALTER TABLE note ADD COLUMN created DATETIME
	`)
	return err
}
