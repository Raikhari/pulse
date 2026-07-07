package main

import (
	"fmt"
	"net/http"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next(w, r)
    }
}

func main() {

	initDB()

	//storage
	http.HandleFunc("/metrics", enableCORS( metricsHandler))
	http.HandleFunc("/metrics/latest", enableCORS(latestHandler))
	http.HandleFunc("/metrics/history", enableCORS(metricsGetHandler))
	http.HandleFunc("/cpu", enableCORS(cpuHandler))
	http.HandleFunc("/stats", enableCORS(statsHandler))

	//API database
	http.HandleFunc("/debug/hosts", enableCORS(debugHostsHandler))
	http.HandleFunc("/debug/latest", enableCORS(debugLatestHandler))
	http.HandleFunc("/debug/dump", enableCORS(debugDumpHandler))

	//Events
	http.HandleFunc("/events", enableCORS(eventsHandler))

	//GET for time series based queries
//	http.HandleFunc("/metrics", metricsQueryHandler)

	fmt.Println("Pulse server listening on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
