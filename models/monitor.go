package models

import (
    "owlhnode/monitor"
)

func GetNodeStats()(data monitor.Monitor) {
	data = monitor.GetLastMonitorInfo()
	return data
}

func AddMonitorFile(anode map[string]string)(err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("MONITOR - AddMonitorFile")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

	err = monitor.AddMonitorFile(anode)
	return err
}

func DeleteMonitorFile(anode map[string]string)(err error) {
	err = monitor.DeleteMonitorFile(anode)
	return err
}

func PingMonitorFiles()(data map[string]map[string]string, err error) {
	data,err = monitor.PingMonitorFiles()
	return data,err
}