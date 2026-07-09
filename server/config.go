package main

type EventConfig struct {
	CPUHighThreshold    float64
	CPUNormalThreshold  float64
	RAMHighThreshold    float64
	RAMNormalThreshold  float64
}

var eventConfig = EventConfig{
	CPUHighThreshold:   80.0,
	CPUNormalThreshold: 70.0,
	RAMHighThreshold:   85.0,
	RAMNormalThreshold: 75.0,
}

type EventConfigResponse struct {
	CPUHighThreshold   float64 `json:"cpu_high_threshold"`
	CPUNormalThreshold float64 `json:"cpu_normal_threshold"`
	RAMHighThreshold   float64 `json:"ram_high_threshold"`
	RAMNormalThreshold float64 `json:"ram_normal_threshold"`
}
