package models

import (
	"owlhnode/ping"
)

func PingService()(err error) {
	err = ping.PingService()
    return err
}

func DeployService(anode map[string]map[string]string)(err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("PING - DeployService")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
	err = ping.DeployService()
    return err
}

func GetMainconfData()(data map[string]map[string]string, err error) {
	data,err = ping.GetMainconfData()
    return data,err
}

func PingPluginsNode()(data map[string]map[string]string ,err error) {
	data, err = ping.PingPluginsNode()
	return data, err
}

func UpdateNodeData(data map[string]map[string]string)(err error) {
    var cc := data
    logs.Info("============")
    logs.Info("PING - UpdateNodeData")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(data,"action")
    delete(data,"controller")
    delete(data,"router")
    
	err = ping.UpdateNodeData(data)
	return err
}