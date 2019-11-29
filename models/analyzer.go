package models

import (
    "owlhnode/analyzer"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"

)

func PingAnalyzer()(data map[string]string ,err error) {
    data, err = analyzer.PingAnalyzer()    
    //changecontrol.ChangeControlInsertData(err, "PingAnalyzer")    
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
    
    changecontrol.ChangeControlInsertData(err, "ChangeAnalyzerStatus")    
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
        
    changecontrol.ChangeControlInsertData(cc, "SyncAnalyzer")    
    return nil
}