package models

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/changeControl"
    "owlhnode/suricata"
)

func GetSuricata() (status map[string]bool, err error) {
    status, err = suricata.Installed()
    // changecontrol.ChangeControlInsertData(err, "GetSuricata")
    return status, err
}

// func GetBPF() (bpf string, err error) {
//     bpf,err = suricata.GetBPF()
//     changecontrol.ChangeControlInsertData(err, "")
// return bpf,err
// }

func SetBPF(anode map[string]string, username string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("SURICATA - SetBPF")
    for key := range cc {
        logs.Info(key + " -> " + cc[key])
    }
    delete(anode, "action")
    delete(anode, "controller")
    delete(anode, "router")

    err = suricata.SetBPF(anode)

    if err != nil {
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    } else {
        cc["actionStatus"] = "success"
    }
    cc["username"] = username
    cc["actionDescription"] = "Set BPF Filter"

    changecontrol.InsertChangeControl(cc)

    // changecontrol.ChangeControlInsertData(err, "SetBPF")
    return err
}

func SyncRulesetFromMaster(n map[string][]byte, username string) (err error) {
    cc := n
    logs.Info("============")
    logs.Info("SURICATA - SyncRulesetFromMaster")
    for key := range cc {
        logs.Info(key + " -> ")
    }
    delete(n, "action")
    delete(n, "controller")
    delete(n, "router")

    err = suricata.SyncRulesetFromMaster(n)

    // if err!=nil {
    //     cc["actionStatus"] = "error"
    //     cc["errorDescription"] = err.Error()
    // }else{
    //     cc["actionStatus"] = "success"
    // }

    // cc["actionDescription"] = "Save new ruleset from master"
    // cc["username"] = username
    // changecontrol.InsertChangeControlByte(cc)
    // changecontrol.ChangeControlInsertData(err, "SyncRulesetFromMaster")
    return err
}

func SaveConfigFile(files map[string][]byte, username string) (err error) {
    cc := files
    logs.Info("============")
    logs.Info("SURICATA - SaveConfigFile")
    for key := range cc {
        logs.Info(key + " -> ")
    }

    err = suricata.SaveConfigFile(files)
    // changecontrol.ChangeControlInsertData(err, "SaveConfigFile")
    return err
}

func RunSuricata(username string) (data string, err error) {
    data, err = suricata.RunSuricata()
    // changecontrol.ChangeControlInsertData(err, "RunSuricata")
    return data, err
}

func SuricataConfigurationTest(uuid string) (data map[string]string, err error) {
    data, err = suricata.SuricataConfigurationTest(uuid)
    // changecontrol.ChangeControlInsertData(err, "RunSuricata")
    return data, err
}

func SuricataDumpCurrentConfig() (data map[string]string, err error) {
    data, err = suricata.SuricataDumpCurrentConfig()
    // changecontrol.ChangeControlInsertData(err, "RunSuricata")
    return data, err
}

func StopSuricata(username string) (data string, err error) {
    data, err = suricata.StopSuricata()
    // changecontrol.ChangeControlInsertData(err, "StopSuricata")
    return data, err
}

func GetSuricataServices(username string) (data map[string]map[string]string, err error) {
    data, err = suricata.GetSuricataServices()
    // changecontrol.ChangeControlInsertData(err, "GetSuricataServices")
    return data, err
}

func StartSuricataMainConf(anode map[string]string, username string) (err error) {
    err = suricata.StartSuricataMainConf(anode)
    changecontrol.ChangeControlInsertData(err, "StartSuricataMainConf")
    return err
}
func StopSuricataMainConf(anode map[string]string, username string) (err error) {
    err = suricata.StopSuricataMainConf(anode)
    changecontrol.ChangeControlInsertData(err, "StopSuricataMainConf")
    return err
}
func KillSuricataMainConf(anode map[string]string, username string) (err error) {
    err = suricata.KillSuricataMainConf(anode)
    changecontrol.ChangeControlInsertData(err, "KillSuricataMainConf")
    return err
}
func ReloadSuricataMainConf(anode map[string]string, username string) (err error) {
    err = suricata.ReloadSuricataMainConf(anode)
    changecontrol.ChangeControlInsertData(err, "ReloadSuricataMainConf")
    return err
}

func GetMD5files(files map[string]map[string]string, username string) (data map[string]map[string]string, err error) {
    data, err = suricata.GetMD5files(files)
    // changecontrol.ChangeControlInsertData(err, "GetMD5files")
    return data, err
}
