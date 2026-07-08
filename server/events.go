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

	highCPU := metrics[0].CPU > eventConfig.CPUHighThreshold
	highRAM := metrics[0].RAM > eventConfig.RAMHighThreshold

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
		if curr.CPU > eventConfig.CPUHighThreshold && !highCPU {
			events = append(events, Event{
				Timestamp: curr.Timestamp,
				Hostname:  curr.Hostname,
				Type:      "high_cpu",
				Message:   "High CPU usage detected",
			})

			highCPU = true
		}

		// Leaving high CPU state
		if curr.CPU < eventConfig.CPUNormalThreshold && highCPU {
			events = append(events, Event{
				Timestamp: curr.Timestamp,
				Hostname:  curr.Hostname,
				Type:      "cpu_normal",
				Message:   "CPU usage returned to normal",
			})
			highCPU = false
		}


		// Entering high RAM state
		if curr.RAM > eventConfig.RAMHighThreshold && !highRAM {
			events = append(events, Event{
				Timestamp: curr.Timestamp,
				Hostname:  curr.Hostname,
				Type:      "high_ram",
				Message:   "High RAM usage detected",
			})
			highRAM = true
		}

		// Leaving high RAM state
		if curr.RAM < eventConfig.RAMNormalThreshold && highRAM {
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

