package models

import (
    "owlhnode/monitor"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func GetNodeStats()(data monitor.Monitor) {
    data = monitor.GetLastMonitorInfo()
    var err error
    err = nil
    //changecontrol.ChangeControlInsertData(err, "GetNodeStats")    
    return data
}

func AddMonitorFile(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("MONITOR - AddMonitorFile")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = monitor.AddMonitorFile(anode)
    //changecontrol.ChangeControlInsertData(err, "AddMonitorFile")    
    return err
}

func DeleteMonitorFile(anode map[string]string)(err error) {
    err = monitor.DeleteMonitorFile(anode)
    //changecontrol.ChangeControlInsertData(err, "DeleteMonitorFile")    
    return err
}

func PingMonitorFiles()(data map[string]map[string]string, err error) {
    data,err = monitor.PingMonitorFiles()
    //changecontrol.ChangeControlInsertData(err, "PingMonitorFiles")    
    return data,err
}