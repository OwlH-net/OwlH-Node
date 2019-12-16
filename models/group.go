package models

import (
    "owlhnode/group"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)

func SyncSuricataGroupValues(data map[string]string) (err error) {
    cc := data
    logs.Info("============")
    logs.Info("GROUP - SyncSuricataGroupValues")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(data,"action")
    delete(data,"controller")
    delete(data,"router")

    logs.Info("Sync Suricata group values")
    err = group.SyncSuricataGroupValues(data)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Sync Suricata Group Values"

    changecontrol.InsertChangeControl(cc)
    return err
}

func SuricataGroupService(data map[string]string) (err error) {
    cc := data
    logs.Info("============")
    logs.Info("GROUP - SuricataGroupService")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(data,"action")
    delete(data,"controller")
    delete(data,"router")

    logs.Info(data["action"]+" Suricata group values")
    err = group.SuricataGroupService(data)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = data["action"]+" Suricata Group Values"

    changecontrol.InsertChangeControl(cc)
    return err
}