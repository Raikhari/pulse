package main

import (
	"encoding/json"
	"net/http"
)

func calculateAverageCPU(data []Metrics) float64 {
	if len(data) == 0 {
		return 0
	}

	var sum float64

	for _, m := range data {
		sum += m.CPU
	}

	return sum / float64(len(data))
}

func calculateAverageRAM(data []Metrics) float64 {
	if len(data) == 0 {
		return 0
	}

	var sum float64

	for _, m := range data {
		sum += m.RAM
	}

	return sum / float64(len(data))
}

func calculateMinRAM(data []Metrics) float64 {

	if len(data) == 0 {
		return 0
	}

	minRAM := data[0].RAM

	for _, m := range data {
		if m.RAM < minRAM {
			minRAM = m.RAM
		}
	}

	return minRAM
}

func calculateMaxRAM(data []Metrics) float64 {

	if len(data) == 0 {
		return 0
	}

	maxRAM := data[0].RAM

	for _, m := range data {
		if m.RAM > maxRAM {
			maxRAM = m.RAM
		}
	}

	return maxRAM
}


func calculateMaxCPU(data []Metrics) float64 {
	if len(data) == 0 {
		return 0
	}

	max := data[0].CPU

	for _, m := range data {
		if m.CPU > max {
			max = m.CPU
		}
	}

	return max
}

func calculateMinCPU(data []Metrics) float64 {

	if len(data) == 0 {
		return 0
	}

	minCPU := data[0].CPU

	for _, m := range data {
		if m.CPU < minCPU {
			minCPU = m.CPU
		}
	}

	return minCPU
}


func calculateAverageLoad1(data []Metrics) float64 {
	if len(data) == 0 {
		return 0
	}

	var sum float64

	for _, m := range data {
		sum += m.Load1
	}

	return sum / float64(len(data))
}

func calculateMaxLoad1(data []Metrics) float64 {
	if len(data) == 0 {
		return 0
	}

	max := data[0].Load1

	for _, m := range data {
		if m.Load1 > max {
			max = m.Load1
		}
	}

	return max
}

func calculateMinLoad1(data []Metrics) float64 {
	if len(data) == 0 {
		return 0
	}

	minLoad1 := data[0].Load1

	for _, m := range data {
		if m.Load1 < minLoad1 {
			minLoad1 = m.Load1
		}
	}
	return minLoad1
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if len(history) == 0 {
		http.Error(w, "no metrics available", http.StatusNotFound)
		return
	}

	host := r.URL.Query().Get("host")

	var data []Metrics

	if host != "" {
		data = history[host]
	} else {
		for _, samples := range history {
			data = append(data, samples...)
		}
	}

	stats := Stats{
		Samples: len(data),

		AvgCPU:  calculateAverageCPU(data),
		MinCPU:  calculateMinCPU(data),
		MaxCPU:  calculateMaxCPU(data),

		AvgRAM:  calculateAverageRAM(data),
		MinRAM: calculateMinRAM(data),
		MaxRAM: calculateMaxRAM(data),

		AvgLoad1: calculateAverageLoad1(data),
		MinLoad1: calculateMinLoad1(data),
		MaxLoad1: calculateMaxLoad1(data),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
