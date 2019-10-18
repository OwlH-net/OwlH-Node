package models

import (
	"owlhnode/ping"
)

func PingService()(err error) {
	err = ping.PingService()
    return err
}

func DeployService()(err error) {
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
	err = ping.UpdateNodeData(data)
	return err
}