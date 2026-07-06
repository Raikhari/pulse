package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"strings"
)

//METRICS
func metricsPostHandler(w http.ResponseWriter, r *http.Request){
	var metrics Metrics

	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}


	//db storage
	_, err = db.Exec(
		`INSERT INTO metrics(hostname, cpu, ram, uptime, load1, timestamp)
		VALUES (?, ?, ?, ?, ?, ?)`,
		metrics.Hostname,
		metrics.CPU,
		metrics.RAM,
		metrics.Uptime,
		metrics.Load1,
		metrics.Timestamp,
	)

	if err != nil {
		fmt.Println("DB insert error:", err)
	}

	fmt.Printf("Received metrics from: %s\n", metrics.Hostname)

	w.WriteHeader(http.StatusOK)

}


func metricsGetHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(`
	SELECT hostname, cpu, ram, uptime, load1, timestamp
	FROM metrics
	ORDER BY timestamp DESC
	LIMIT 100
	`)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	defer rows.Close()

	var result []Metrics

	for rows.Next() {
		var m Metrics
		rows.Scan(&m.Hostname, &m.CPU, &m.RAM, &m.Uptime, &m.Load1, &m.Timestamp)
		result = append(result, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		metricsPostHandler(w, r)

	case http.MethodGet:
		metricsQueryHandler(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}


func latestHandler(w http.ResponseWriter, r *http.Request){

	host := r.URL.Query().Get("host")

	if host == "" {
		http.Error(w, "host required", http.StatusBadRequest)
		return
	}

	var m Metrics

	err := db.QueryRow(`
	SELECT hostname, cpu, ram, uptime, load1, timestamp
	FROM metrics
	WHERE hostname = ?
	ORDER BY timestamp DESC
	LIMIT 1
	`, host).Scan(
		&m.Hostname,
		&m.CPU,
		&m.RAM,
		&m.Uptime,
		&m.Load1,
		&m.Timestamp,
	)

	if err != nil {
		http.Error(w, "no data available", http.StatusNotFound)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

//cpu handler
func cpuHandler(w http.ResponseWriter, r *http.Request) {

	host := r.URL.Query().Get("host")

	if host == "" {
		http.Error(w, "host required", http.StatusBadRequest)
		return
	}

	type point struct {
		Timestamp int64   `json:"timestamp"`
		CPU       float64 `json:"cpu"`
	}


	rows, err := db.Query(`
	SELECT timestamp, cpu
	FROM metrics
	WHERE hostname = ?
	ORDER BY timestamp ASC
	`, host)

	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	defer rows.Close()

	var series []point

	for rows.Next() {
		var p point
		rows.Scan(&p.Timestamp, &p.CPU)
		series = append(series, p)
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}

//DATABASE
//debug handlers
func debugHostsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
	SELECT DISTINCT hostname
	FROM metrics
	`)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var hosts []string

	for rows.Next() {
		var h string
		rows.Scan(&h)
		hosts = append(hosts, h)
	}

	json.NewEncoder(w).Encode(hosts)
}

func debugLatestHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	if host == "" {
		http.Error(w, "missing host", 400)
		return
	}

	row := db.QueryRow(`
	SELECT hostname, cpu, ram, uptime, load1, timestamp
	FROM metrics
	WHERE hostname = ?
	ORDER BY timestamp DESC
	LIMIT 1
	`, host)

	var m Metrics
	err := row.Scan(
		&m.Hostname,
		&m.CPU,
		&m.RAM,
		&m.Uptime,
		&m.Load1,
		&m.Timestamp,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(m)
}

func debugDumpHandler(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "100"
	}

	rows, err := db.Query(`
	SELECT hostname, cpu, ram, uptime, load1, timestamp
	FROM metrics
	ORDER BY timestamp DESC
	LIMIT ?
	`, limit)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var out []Metrics

	for rows.Next() {
		var m Metrics
		rows.Scan(&m.Hostname, &m.CPU, &m.RAM, &m.Uptime, &m.Load1, &m.Timestamp)
		out = append(out, m)
	}

	json.NewEncoder(w).Encode(out)
}

func metricsQueryHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	limit := r.URL.Query().Get("limit")
	hours := r.URL.Query().Get("hours")

	if limit == "" {
		limit = "100"
	}

	query := `
	SELECT hostname, cpu, ram, uptime, load1, timestamp
	FROM metrics
	`

	var args []any
	var conditions []string

	if host != "" {
		conditions = append(conditions, "hostname = ?")
		args = append(args, host)
	}

	if hours != "" {
		h, err := strconv.Atoi(hours)
		if err == nil {
			since := time.Now().Unix() - int64(h*3600)
			conditions = append(conditions, "timestamp >= ?")
			args = append(args, since)
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY timestamp DESC LIMIT ?"
	args = append(args, limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	defer rows.Close()

	var result []Metrics

	for rows.Next() {
		var m Metrics
		rows.Scan(&m.Hostname, &m.CPU, &m.RAM, &m.Uptime, &m.Load1, &m.Timestamp)
		result = append(result, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
