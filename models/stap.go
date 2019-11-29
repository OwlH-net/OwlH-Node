package models

import (
    "owlhnode/stap"
    "owlhnode/changeControl"
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
    // changecontrol.ChangeControlInsertData(err, "AddServer")    
    return err
}
func GetAllServers()(servers map[string]map[string]string, err error) {
    logs.Info("GetAllServers stap into Node file")
    servers,err = stap.GetAllServers()  
    // changecontrol.ChangeControlInsertData(err, "GetAllServers")    
    return servers,err
}

func GetServer(serveruuid string)(servers map[string]map[string]string, err error) {
    logs.Info("GetAllServers stap into Node file")
    servers,err = stap.GetServer(serveruuid)
    // changecontrol.ChangeControlInsertData(err, "GetServer")    
    return servers,err
}

func PingStap(uuid string) (status map[string]bool, err error) {
    status, err = stap.PingStap(uuid)
    // changecontrol.ChangeControlInsertData(err, "PingStap")    
    return status, err
}

func PingServerStap(server string) (status map[string]string, err error) {
    status, err = stap.PingServerStap(server)
    // changecontrol.ChangeControlInsertData(err, "PingServerStap")    
    return status, err
}

func RunStap(uuid string) (data string, err error) {
    logs.Info("Run RunStap system into node server")
    data,err = stap.RunStap(uuid)
    // changecontrol.ChangeControlInsertData(err, "RunStap")    
    return data,err
}

func StopStap(uuid string) (data string, err error) {
    logs.Info("Stops StopStap system into node server")
    data,err = stap.StopStap(uuid)
    // changecontrol.ChangeControlInsertData(err, "StopStap")    
    return data,err
}

func RunStapServer(serveruuid string) (data string, err error) {
    logs.Info("Run RunStapServer system into node server")
    data,err = stap.RunStapServer(serveruuid)
    // changecontrol.ChangeControlInsertData(err, "RunStapServer")    
    return data,err
}

func DeleteStapServer(serveruuid string) (data string, err error) {
    data,err = stap.DeleteStapServer(serveruuid)
    // changecontrol.ChangeControlInsertData(err, "DeleteStapServer")    
    return data,err
}

func StopStapServer(serveruuid string) (data string, err error) {
    logs.Info("Stops StopStapServer system into node server")
    data,err = stap.StopStapServer(serveruuid)
    // changecontrol.ChangeControlInsertData(err, "StopStapServer")    
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
    // changecontrol.ChangeControlInsertData(err, "EditStapServer")    
    return err
}