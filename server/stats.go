package main

import (
    "database/sql"
    "encoding/json"
    "net/http"
)


func statsHandler(w http.ResponseWriter, r *http.Request) {

    host := r.URL.Query().Get("host")

    query := `
        SELECT
            COUNT(*),

            AVG(cpu),
            MIN(cpu),
            MAX(cpu),

            AVG(ram),
            MIN(ram),
            MAX(ram),

            AVG(load1),
            MIN(load1),
            MAX(load1)
        FROM metrics
    `

    var row *sql.Row

    if host != "" {
        query += " WHERE hostname = ?"
        row = db.QueryRow(query, host)
    } else {
        row = db.QueryRow(query)
    }

    var stats Stats

    err := row.Scan(
        &stats.Samples,

        &stats.AvgCPU,
        &stats.MinCPU,
        &stats.MaxCPU,

        &stats.AvgRAM,
        &stats.MinRAM,
        &stats.MaxRAM,

        &stats.AvgLoad1,
        &stats.MinLoad1,
        &stats.MaxLoad1,
    )

    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}


