package models

import (
    "owlhnode/monitor"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func GetNodeStats()(data monitor.Monitor) {
    data = monitor.GetLastMonitorInfo()
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
    
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Start monitoring a file"

    changecontrol.InsertChangeControl(cc)

    //changecontrol.ChangeControlInsertData(err, "AddMonitorFile")    
    return err
}

func DeleteMonitorFile(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("MONITOR - DeleteMonitorFile")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    
    err = monitor.DeleteMonitorFile(anode)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Delete traffic transport values"

    changecontrol.InsertChangeControl(cc)

    //changecontrol.ChangeControlInsertData(err, "DeleteMonitorFile")    
    return err
}

func PingMonitorFiles()(data map[string]map[string]string, err error) {
    data,err = monitor.PingMonitorFiles()
    //changecontrol.ChangeControlInsertData(err, "PingMonitorFiles")    
    return data,err
}