package models

import (
    "owlhnode/wazuh"
    "github.com/astaxie/beego/logs"
)


func GetWazuh() (status map[string]bool, err error) {
	return wazuh.Installed()
}

func RunWazuh() (data string, err error) {
    logs.Info("Run RunWazuh system into node server")
    data,err = wazuh.RunWazuh()
    return data,err
}

func StopWazuh() (data string, err error) {
    logs.Info("Stops StopWazuh system into node server")
    data,err = wazuh.StopWazuh()
    return data,err
}

func PingWazuhFiles() (files map[string]string, err error) {
    files, err = wazuh.PingWazuhFiles()
    return files ,err
}

func DeleteWazuhFile(file map[string]interface{})(err error) {
    err = wazuh.DeleteWazuhFile(file)
    return err
}