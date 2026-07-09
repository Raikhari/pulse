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

	createConfigTable := `
	CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value REAL
	);
	`

	_, err = db.Exec(createConfigTable)
	if err != nil {
		log.Fatal(err)
	}

	loadEventConfig()
}

func loadEventConfig() {
	rows, err := db.Query(`
	SELECT key, value
	FROM config
	`)
	if err != nil {
		log.Println("Config load error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var value float64

		err := rows.Scan(&key, &value)
		if err != nil {
			continue
		}

		switch key {
		case "cpu_high_threshold":
			eventConfig.CPUHighThreshold = value
		case "cpu_normal_threshold":
			eventConfig.CPUNormalThreshold = value
		case "ram_high_threshold":
			eventConfig.RAMHighThreshold = value
		case "ram_normal_threshold":
			eventConfig.RAMNormalThreshold = value
		}
	}
}

func saveEventConfig() error {
	values := map[string]float64{
		"cpu_high_threshold":    eventConfig.CPUHighThreshold,
		"cpu_normal_threshold": eventConfig.CPUNormalThreshold,
		"ram_high_threshold":   eventConfig.RAMHighThreshold,
		"ram_normal_threshold": eventConfig.RAMNormalThreshold,
	}

	for key, value := range values {
		_, err := db.Exec(`
			INSERT OR REPLACE INTO config(key, value)
			VALUES (?, ?)
		`, key, value)

		if err != nil {
			return err
		}
	}

	return nil
}
