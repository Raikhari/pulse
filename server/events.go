package main

type Event struct {
	Timestamp int64  `json:"timestamp"`
	Hostname  string `json:"hostname"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

func GenerateEvents(metrics []Metrics) []Event {

	var events []Event

	if len(metrics) < 2 {
		return events
	}

	cpuHighThreshold := 80.0
	cpuNormalThreshold := 70.0

	highCPU := metrics[0].CPU > cpuHighThreshold

	ramHighThreshold := 85.0
	ramNormalThreshold := 75.0

	highRAM := metrics[0].RAM > ramHighThreshold

	for i := 1; i < len(metrics); i++ {

		prev := metrics[i-1]
		curr := metrics[i]

		// reboot detection
		if curr.Uptime < prev.Uptime {
			events = append(events, Event{
				Timestamp: curr.Timestamp,
				Hostname:  curr.Hostname,
				Type:      "reboot",
				Message:   "System reboot detected",
			})
		}

		// Entering high CPU state
		if curr.CPU > cpuHighThreshold && !highCPU {
			events = append(events, Event{
				Timestamp: curr.Timestamp,
				Hostname:  curr.Hostname,
				Type:      "high_cpu",
				Message:   "High CPU usage detected",
			})

			highCPU = true
		}

		// Leaving high CPU state
		if curr.CPU < cpuNormalThreshold && highCPU {
			events = append(events, Event{
				Timestamp: curr.Timestamp,
				Hostname:  curr.Hostname,
				Type:      "cpu_normal",
				Message:   "CPU usage returned to normal",
			})
			highCPU = false
		}


		// Entering high RAM state
		if curr.RAM > ramHighThreshold && !highRAM {
			events = append(events, Event{
				Timestamp: curr.Timestamp,
				Hostname:  curr.Hostname,
				Type:      "high_ram",
				Message:   "High RAM usage detected",
			})
			highRAM = true
		}

		// Leaving high CPU state
		if curr.RAM < ramNormalThreshold && highRAM {
			events = append(events, Event{
				Timestamp: curr.Timestamp,
				Hostname:  curr.Hostname,
				Type:      "ram_normal",
				Message:   "RAM usage returned to normal",
			})
			highRAM = false
		}
	}
	return events
}

