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

func LoadConfig(anode map[string]string) () {
    cc := anode
    logs.Info("============")
    logs.Info("hwaddmng - LoadConfig")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }

    logs.Info("MODELS -> Load default config")
    pcap.LoadConfig()

    cc["actionStatus"] = "success"
    cc["actionDescription"] = "Load default config"

    changecontrol.InsertChangeControl(cc)
}

func Config(anodeIface map[string]interface{}, anode map[string]string) () {
    cc := anode
    logs.Info("============")
    logs.Info("hwaddmng - Config")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }

    logs.Info("MODELS -> Load default config")
    pcap.Config(anodeIface)

    cc["actionStatus"] = "success"
    cc["actionDescription"] = "Load default config"

    changecontrol.InsertChangeControl(cc)
}

func Db(anodeIface map[string]string, anode map[string]string) () {
    cc := anode
    logs.Info("============")
    logs.Info("hwaddmng - Config")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }

    logs.Info("MODELS -> Db management")
    pcap.Db(anodeIface)

    cc["actionStatus"] = "success"
    cc["actionDescription"] = "Db management"

    changecontrol.InsertChangeControl(cc)
}