package models

import (
	"owlhnode/net"
//    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func GetNetworkData()(values map[string]string, err error) {
	values,err = net.GetNetworkData()
    return values,err
}

func LoadNetworkValuesSelected()(values map[string]map[string]string, err error) {
	values,err = net.LoadNetworkValuesSelected()
    return values,err
}

func UpdateNetworkInterface(data map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("NET - UpdateNetworkInterface")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

	err = net.UpdateNetworkInterface(data)
    return err
}

// func LoadNetworkValuesSuricata()(values map[string]map[string]string, err error) {
// 	values,err = net.LoadNetworkValuesSuricata()
//     return values,err
// }