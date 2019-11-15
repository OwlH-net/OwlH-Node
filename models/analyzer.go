package models

import (
    "owlhnode/analyzer"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"

)

func PingAnalyzer()(data map[string]string ,err error) {
    data, err = analyzer.PingAnalyzer()    
    return data, err
}

func ChangeAnalyzerStatus(uuid map[string]string) (err error) {
    logs.Info("============")
    logs.Info("ANALYZER - ChangeAnalyzerStatus")
    cc := uuid
    for key :=range uuid {
        logs.Info(key +" -> "+cc[key])
    }
    
    err = analyzer.ChangeAnalyzerStatus(uuid)
    
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["actionDescription"] = "Change Analyzer Status Enable/Disable"
    
    controlError := changecontrol.InsertChangeControl(cc)
    if controlError!=nil { logs.Error("AddPluginService controlError: "+controlError.Error()) }
    
    return err
}

func SyncAnalyzer(file map[string][]byte) (err error) {
    cc := make(map[string]string)
    cc["action"] = "PUT"
    cc["controller"] = "ANALYZER"
    cc["router"] = "@router /SyncAnalyzer [put]"
    logs.Info("============")
    logs.Info("ANALYZER - SyncAnalyzer")
    logs.Info("file - conf/analyzer.json")
    
    err = analyzer.SyncAnalyzer(file)
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["actionDescription"] = "Sync analyzer file"
    
    controlError := changecontrol.InsertChangeControl(cc)
    if controlError!=nil { logs.Error("SyncAnalyzer controlError: "+controlError.Error()) }
    
    return err
}