package main

type Metrics struct {
    Hostname  string  `json:"hostname"`
    CPU       float64 `json:"cpu"`
    RAM       float64 `json:"ram"`
    Uptime    float64 `json:"uptime"`
    Load1     float64 `json:"load1"`
    Timestamp int64   `json:"timestamp"`
}

type Stats struct {
    Samples int     `json:"samples"`
    AvgCPU  float64 `json:"avg_cpu"`
    MinCPU  float64 `json:"min_cpu"`
    MaxCPU  float64 `json:"max_cpu"`
    
    AvgRAM  float64 `json:"avg_ram"`
    MinRAM  float64 `json:"min_ram"`
    MaxRAM  float64 `json:"max_ram"`

    AvgLoad1 float64 `json:"avg_load1"`
    MinLoad1 float64 `json:"min_load1"`
    MaxLoad1 float64 `json:"max_load1"`

}
