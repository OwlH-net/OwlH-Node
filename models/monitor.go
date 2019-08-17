package models

import (
    "owlhnode/monitor"
)

func GetNodeStats()(data monitor.Monitor) {
	data = monitor.GetLastMonitorInfo()
	return data
}

