package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error

	db, err = sql.Open("sqlite3", "./metrics.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hostname TEXT,
		cpu REAL,
		ram REAL,
		uptime REAL,
		load1 REAL,
		timestamp INTEGER
	);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}
