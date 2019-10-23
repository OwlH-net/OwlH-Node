package models

import (
    "owlhnode/analyzer"
//    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"

)

func PingAnalyzer()(data map[string]string ,err error) {
	data, err = analyzer.PingAnalyzer()	
	return data, err
}

func ChangeAnalyzerStatus(uuid map[string]string) (err error) {
    
    logs.Info("============")
    logs.Info("ANALYZER - ChangeAnalyzerStatus")
    for key :=range uuid {
        logs.Info(key +" -> "+uuid[key])
    }
	err = analyzer.ChangeAnalyzerStatus(uuid)
	return err
}