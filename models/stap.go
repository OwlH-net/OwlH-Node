package models

import (
    "owlhnode/stap"
//    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func AddServer(anode map[string]string) (err error) {
    logs.Info("AddServer stap into Node file")
    cc := anode
    logs.Info("============")
    logs.Info("STAP - AddServer")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = stap.AddServer(anode)
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

func EditStapServer(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("STAP - EditStapServer")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = stap.EditStapServer(anode)
    return err
}