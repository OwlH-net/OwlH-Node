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