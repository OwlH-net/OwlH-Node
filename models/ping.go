package models

import (
    "owlhnode/ping"
    // "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func PingService()(err error) {
    err = ping.PingService()
    //changecontrol.ChangeControlInsertData(err, "PingService")    
    return err
}

func DeployService(anode map[string]map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PING - DeployService")
    for key :=range cc {
        logs.Info(key +" -> ")
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    err = ping.DeployService()

    // if err!=nil { 
    //     cc["actionStatus"] = "error"
    //     cc["errorDescription"] = err.Error()
    // }else{
    //     cc["actionStatus"] = "success"
    // }

    // cc["actionDescription"] = "Deploy Service"

    // changecontrol.InsertChangeControlExtra(cc)
    //changecontrol.ChangeControlInsertData(err, "DeployService")    
    return err
}

func GetMainconfData()(data map[string]map[string]string, err error) {
    data,err = ping.GetMainconfData()
    //changecontrol.ChangeControlInsertData(err, "GetMainconfData")    
    return data,err
}

func PingPluginsNode()(data map[string]map[string]string ,err error) {
    data, err = ping.PingPluginsNode()
    //changecontrol.ChangeControlInsertData(err, "PingPluginsNode")    
    return data, err
}

func UpdateNodeData(data map[string]map[string]string)(err error) {
    cc := data
    logs.Info("============")
    logs.Info("PING - UpdateNodeData")
    for key :=range cc {
        logs.Info(key +" -> ")
    }
    delete(data,"action")
    delete(data,"controller")
    delete(data,"router")
    
    err = ping.UpdateNodeData(data)

    // if err!=nil { 
    //     cc["actionStatus"] = "error"
    //     cc["errorDescription"] = err.Error()
    // }else{
    //     cc["actionStatus"] = "success"
    // }

    // cc["actionDescription"] = "update node info"

    // //complex map
    // changecontrol.InsertChangeControlExtra(cc)
    //changecontrol.ChangeControlInsertData(err, "UpdateNodeData")    
    return err
}

func SaveNodeInformation(anode map[string]map[string]string)(err error) {
    err = ping.SaveNodeInformation(anode)
    //changecontrol.ChangeControlInsertData(err, "SaveNodeInformation")    
    return err
}

func DeleteNode()(err error) {
    err = ping.DeleteNode()
    //changecontrol.ChangeControlInsertData(err, "DeleteNode")    
    return err
}