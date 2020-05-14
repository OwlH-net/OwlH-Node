package models

import (
    "owlhnode/analyzer"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)

func PingAnalyzer(username string)(data map[string]string, err error) {
    data, err = analyzer.PingAnalyzer()
    //changecontrol.ChangeControlInsertData(err, "PingAnalyzer", username)    
    return data, err
}

func ChangeAnalyzerStatus(anode map[string]string, username string) (err error) {
    logs.Info("============")
    logs.Info("ANALYZER - ChangeAnalyzerStatus")
    cc := anode
    for key :=range anode {
        logs.Info(key +" -> "+cc[key])
    }
    
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    
    err = analyzer.ChangeAnalyzerStatus(anode)
    
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["username"] = username
    cc["actionDescription"] = "Change Analyzer Status"


    changecontrol.InsertChangeControl(cc)
    return err
}

func SyncAnalyzer(file map[string][]byte, username string) (err error) {
    cc := make(map[string]string)
    cc["action"] = "PUT"
    cc["controller"] = "ANALYZER"
    cc["router"] = "@router /SyncAnalyzer [put]"
    logs.Info("============")
    logs.Info("ANALYZER - SyncAnalyzer")
    logs.Info("file - conf/analyzer.json")
    //TODO action
    err = analyzer.SyncAnalyzer(file)
    
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["username"] = username
    cc["actionDescription"] = "sync Analyzer configuration"

    changecontrol.InsertChangeControl(cc)    
    return nil
}