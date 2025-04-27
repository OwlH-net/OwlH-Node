package models

import (
	changecontrol "github.com/OwlH-net/OwlH-Node/changeControl"
	"github.com/OwlH-net/OwlH-Node/net"
	"github.com/astaxie/beego/logs"
)

func GetNetworkData(username string) (values map[string]string, err error) {
	values, err = net.GetNetworkData()
	//changecontrol.ChangeControlInsertData(err, "GetNetworkData")
	return values, err
}

func LoadNetworkValuesSelected(username string) (values map[string]map[string]string, err error) {
	values, err = net.LoadNetworkValuesSelected()
	//changecontrol.ChangeControlInsertData(err, "LoadNetworkValuesSelected")
	return values, err
}

func UpdateNetworkInterface(data map[string]string, username string) (err error) {
	cc := data
	logs.Info("============")
	logs.Info("NET - UpdateNetworkInterface")
	for key := range cc {
		logs.Info(key + " -> " + cc[key])
	}
	delete(data, "action")
	delete(data, "controller")
	delete(data, "router")

	err = net.UpdateNetworkInterface(data)

	if err != nil {
		cc["actionStatus"] = "error"
		cc["errorDescription"] = err.Error()
	} else {
		cc["actionStatus"] = "success"
	}
	cc["username"] = username
	cc["actionDescription"] = "Update default network Interface"

	changecontrol.InsertChangeControl(cc)
	//changecontrol.ChangeControlInsertData(err, "UpdateNetworkInterface")
	return err
}

// func LoadNetworkValuesSuricata()(values map[string]map[string]string, err error) {
//     values,err = net.LoadNetworkValuesSuricata()
//     changecontrol.ChangeControlInsertData(err, "")
// return values,err
// }
