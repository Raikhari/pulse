package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Metrics struct {
	Hostname string `json:"hostname"`
	CPU      float64 `json:"cpu"`
	RAM      float64 `json:"ram"`
    	Uptime   float64 `json:"uptime"`
	Load1	 float64 `json:"load1"`
	Timestamp int64  `json:"timestamp"`
}

type Stats struct {
	Samples int `json:"samples"`
	AvgCPU float64 `json:"avg_cpu"`
	MaxCPU float64 `json:"max_cpu"`
	AvgRAM float64 `json:"avg_ram"`
}

var history []Metrics

//METRICS
func metricsPostHandler(w http.ResponseWriter, r *http.Request){
	var metrics Metrics

	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	history = append(history, metrics)

	fmt.Printf("Received metrics from: %s\n", metrics.Hostname)

	w.WriteHeader(http.StatusOK)

}


func metricsGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		metricsPostHandler(w, r)

	case http.MethodGet:
		metricsGetHandler(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func historyHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func latestHandler(w http.ResponseWriter, r *http.Request){
	
	if len(history) == 0 {
		http.Error(w, "no metrics available", http.StatusNotFound)
		return
	}

	latest := history[len(history)-1]
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(latest)
}

//cpu handler
func cpuHandler(w http.ResponseWriter, r *http.Request) {
    var series []struct {
        Timestamp int64   `json:"timestamp"`
        CPU       float64 `json:"cpu"`
    }

    for _, m := range history {
        series = append(series, struct {
            Timestamp int64   `json:"timestamp"`
            CPU       float64 `json:"cpu"`
        }{
            Timestamp: m.Timestamp,
            CPU:       m.CPU,
        })
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(series)
}

//STATS
func statsHandler(w http.ResponseWriter, r *http.Request) {

	if len(history) == 0 {
		http.Error(w, "no metrics available", http.StatusNotFound)
		return
	}

	var sumCPU float64
	var sumRAM float64
	var maxCPU float64

	for _, m := range history {

		sumCPU += m.CPU
		sumRAM += m.RAM

		if m.CPU > maxCPU {
			maxCPU = m.CPU
		}
	}

	stats := Stats{
		Samples: len(history),
		AvgCPU:  sumCPU / float64(len(history)),
		MaxCPU:  maxCPU,
		AvgRAM:  sumRAM / float64(len(history)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}


func main() {
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/metrics/latest", latestHandler)
	http.HandleFunc("/metrics/history", historyHandler)
	http.HandleFunc("/cpu", cpuHandler)

	http.HandleFunc("/stats", statsHandler)
	
	fmt.Println("Pulse server listening on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

