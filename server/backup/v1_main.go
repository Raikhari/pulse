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
}

var latest Metrics

func metricsPostHandler(w http.ResponseWriter, r *http.Request){
	var metrics Metrics

	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	latest = metrics

	fmt.Printf("Received metrics from: %s\n", metrics.Hostname)

	w.WriteHeader(http.StatusOK)

}


func metricsGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(latest)
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

func main() {
	http.HandleFunc("/metrics", metricsHandler)

	fmt.Println("Pulse server listening on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

