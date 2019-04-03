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
func GetAllServers()(servers *map[string]map[string]string, err error) {
    logs.Info("GetAllServers stap into Node file")
	servers,err = stap.GetAllServers()
    return servers,err
}

func GetServer(serveruuid string)(servers *map[string]map[string]string, err error) {
    logs.Info("GetAllServers stap into Node file")
	servers,err = stap.GetServer(serveruuid)
    return servers,err
}

func PingStap(uuid string) (status map[string]bool, err error) {
	status, err = stap.PingStap(uuid)
	return status, err
}

func PingServerStap(server string) (status map[string]bool, err error) {
	status, err = stap.PingServerStap(server)
	return status, err
}

func RunStap(uuid string) (data string, err error) {
    logs.Info("Run RunStap system into node server")
    data,err = stap.RunStap(uuid)
    return data,err
}

func StopStap(uuid string) (data string, err error) {
    logs.Info("Stops StopStap system into node server")
    data,err = stap.StopStap(uuid)
    return data,err
}

func RunStapServer(serveruuid string) (data string, err error) {
    logs.Info("Run RunStapServer system into node server")
    data,err = stap.RunStapServer(serveruuid)
    return data,err
}

func StopStapServer(serveruuid string) (data string, err error) {
    logs.Info("Stops StopStapServer system into node server")
    data,err = stap.StopStapServer(serveruuid)
    return data,err
}