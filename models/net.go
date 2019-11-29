package models

import (
    "owlhnode/net"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func GetNetworkData()(values map[string]string, err error) {
    values,err = net.GetNetworkData()
    //changecontrol.ChangeControlInsertData(err, "GetNetworkData")    
    return values,err
}

func LoadNetworkValuesSelected()(values map[string]map[string]string, err error) {
    values,err = net.LoadNetworkValuesSelected()
    //changecontrol.ChangeControlInsertData(err, "LoadNetworkValuesSelected")    
    return values,err
}

func UpdateNetworkInterface(data map[string]string) (err error) {
    cc := data
    logs.Info("============")
    logs.Info("NET - UpdateNetworkInterface")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(data,"action")
    delete(data,"controller")
    delete(data,"router")

    err = net.UpdateNetworkInterface(data)
    //changecontrol.ChangeControlInsertData(err, "UpdateNetworkInterface")    
    return err
}

// func LoadNetworkValuesSuricata()(values map[string]map[string]string, err error) {
//     values,err = net.LoadNetworkValuesSuricata()
//     changecontrol.ChangeControlInsertData(err, "")    
    // return values,err
// }