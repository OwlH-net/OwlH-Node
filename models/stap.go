package models

import (
	"owlhnode/stap"
	"github.com/astaxie/beego/logs"
)

func AddServer(elem map[string]string) (err error) {
    logs.Info("AddServer stap into Node file")
	err = stap.AddServer(elem)
    return err
}
func GetAllServers()(servers map[string]map[string]string, err error) {
    logs.Info("GetAllServers stap into Node file")
	servers,err = stap.GetAllServers()
	if err != nil {
		return nil, err
	}
    return servers,err
}

func GetServer(serveruuid string)(servers map[string]map[string]string, err error) {
    logs.Info("GetAllServers stap into Node file")
	servers,err = stap.GetServer(serveruuid)
	if err != nil {
		return nil, err
	}
    return servers,err
}

func PingStap(uuid string) (status map[string]bool, err error) {
	status, err = stap.PingStap(uuid)
	if err != nil {
		return nil, err
	}
	return status, err
}

func PingServerStap(server string) (status map[string]string, err error) {
	status, err = stap.PingServerStap(server)
	if err != nil {
		return nil, err
	}
	return status, err
}

func RunStap(uuid string) (data string, err error) {
    logs.Info("Run RunStap system into node server")
	data,err = stap.RunStap(uuid)
	if err != nil {
		return "", err
	}
    return data,err
}

func StopStap(uuid string) (data string, err error) {
    logs.Info("Stops StopStap system into node server")
	data,err = stap.StopStap(uuid)
	if err != nil {
		return "", err
	}
    return data,err
}

func RunStapServer(serveruuid string) (data string, err error) {
    logs.Info("Run RunStapServer system into node server")
	data,err = stap.RunStapServer(serveruuid)
	if err != nil {
		return "", err
	}
    return data,err
}

func DeleteStapServer(serveruuid string) (data string, err error) {
	data,err = stap.DeleteStapServer(serveruuid)
	if err != nil {
		return "", err
	}
    return data,err
}

func StopStapServer(serveruuid string) (data string, err error) {
    logs.Info("Stops StopStapServer system into node server")
    data,err = stap.StopStapServer(serveruuid)
	if err != nil {
		return "", err
	}
	return data,err
}

func EditStapServer(data map[string]string) (err error) {
    err = stap.EditStapServer(data)
    return err
}