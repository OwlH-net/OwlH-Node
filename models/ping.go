package models

import (
	"github.com/OwlH-net/OwlH-Node/ping"
	// "github.com/OwlH-net/OwlH-Node/changeControl"
	"github.com/astaxie/beego/logs"
)

func PingService(username string) (err error) {
	err = ping.PingService()
	//changecontrol.ChangeControlInsertData(err, "PingService")
	return err
}

func DeployService(anode map[string]map[string]string, username string) (err error) {
	cc := anode
	logs.Info("============")
	logs.Info("PING - DeployService")
	for key := range cc {
		logs.Info(key + " -> ")
	}
	delete(anode, "action")
	delete(anode, "controller")
	delete(anode, "router")
	err = ping.DeployService()

	// if err!=nil {
	//     cc["actionStatus"] = "error"
	//     cc["errorDescription"] = err.Error()
	// }else{
	//     cc["actionStatus"] = "success"
	// }
	// cc["username"] = username
	// cc["actionDescription"] = "Deploy Service"

	// changecontrol.InsertChangeControlExtra(cc)
	//changecontrol.ChangeControlInsertData(err, "DeployService")
	return err
}

func GetMainconfData(username string) (data map[string]map[string]string, err error) {
	data, err = ping.GetMainconfData()
	//changecontrol.ChangeControlInsertData(err, "GetMainconfData")
	return data, err
}

func PingPluginsNode(username string) (data map[string]map[string]string, err error) {
	data, err = ping.PingPluginsNode()
	//changecontrol.ChangeControlInsertData(err, "PingPluginsNode")
	return data, err
}

func UpdateNodeData(data map[string]map[string]string, username string) (err error) {
	cc := data
	logs.Info("============")
	logs.Info("PING - UpdateNodeData")
	for key := range cc {
		logs.Info(key + " -> ")
	}
	delete(data, "action")
	delete(data, "controller")
	delete(data, "router")

	err = ping.UpdateNodeData(data)

	// if err!=nil {
	//     cc["actionStatus"] = "error"
	//     cc["errorDescription"] = err.Error()
	// }else{
	//     cc["actionStatus"] = "success"
	// }
	// cc["username"] = username
	// cc["actionDescription"] = "update node info"

	// //complex map
	// changecontrol.InsertChangeControlExtra(cc)
	//changecontrol.ChangeControlInsertData(err, "UpdateNodeData")
	return err
}

func SaveNodeInformation(anode map[string]map[string]string, username string) (err error) {
	err = ping.SaveNodeInformation(anode)
	//changecontrol.ChangeControlInsertData(err, "SaveNodeInformation")
	return err
}

func DeleteNode(masterID string, username string) (err error) {
	err = ping.DeleteNode(masterID)
	//changecontrol.ChangeControlInsertData(err, "DeleteNode")
	return err
}
