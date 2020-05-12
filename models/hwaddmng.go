package models

import (
    "owlhnode/pcap"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)

func AddMacIp(data map[string]string) (err error) {
    cc := data
    logs.Info("============")
    logs.Info("hwaddmng - AddMacIp")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(data,"action")
    delete(data,"controller")
    delete(data,"router")

    logs.Info("MODELS -> Add IP and MAC to Knownmacs")
    err = pcap.AddMacIp(data)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Add IP and MAC to Knownmacs"

    changecontrol.InsertChangeControl(cc)
    return err
}