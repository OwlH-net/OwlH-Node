package models

import (
    "owlhnode/plugin"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)

func ChangeServiceStatus(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - ChangeServiceStatus")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = plugin.ChangeServiceStatus(anode)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Change " + cc["service"] + " status"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "ChangeServiceStatus")    
    return err
}

// curl -X PUT \
//   https://52.47.197.22:50002/node/plugin/ChangeMainServiceStatus \
//   -H 'Content-Type: application/json' \
//   -d '{
//     "service": "suricata",
//     "param": "status",
//     "value": "enabled"
// }
func ChangeMainServiceStatus(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - ChangeMainServiceStatus")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = plugin.ChangeMainServiceStatus(anode)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "change " + cc["service"] + " status to " + cc["value"]

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "ChangeMainServiceStatus")    
    return err
}

func DeleteService(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - DeleteService")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = plugin.DeleteService(anode)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Delete Service"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "DeleteService")    
    return err
}

func AddPluginService(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - AddPluginService")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = plugin.AddPluginService(anode)
    
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Add Service"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "AddPluginService")    
    return err
}

func SaveSuricataInterface(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - SaveSuricataInterface")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    err = plugin.SaveSuricataInterface(anode)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Update Suricata Interface"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "SaveSuricataInterface")    
    return err
}

func DeployStapService(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - DeployStapService")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = plugin.DeployStapService(anode)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Start Software TAP service"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "DeployStapService")    
    return err
}

func StopStapService(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - StopStapService")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = plugin.StopStapService(anode)
    
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Stop Software TAP Service"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "StopStapService")    
    return err
}

func ModifyStapValues(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - ModifyStapValues")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = plugin.ModifyStapValues(anode)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Modify Software TAP configuration"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "ModifyStapValues")    
    return err
}

// curl -X PUT \
//   https://52.47.197.22:50002/node/plugin/changeSuricataTable \
//   -H 'Content-Type: application/json' \
//   -d '{
//     "uuid": "suricata",
//     "status": "none"
// }
func ChangeSuricataTable(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - ChangeSuricataTable")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = plugin.ChangeSuricataTable(anode)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Change Suricata Table"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "ChangeSuricataTable")    
    return err
}

// curl -X PUT \
//   https://52.47.197.22:50002/node/plugin/getCommands \
//   -H 'Content-Type: application/json' \
//   -d '{
//     "uuid": "suricata",
//     "service": "service"
// }
func GetServiceCommands(anode map[string]string)(data map[string]map[string]string, err error) {
    cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - GetServiceCommands")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    data,err = plugin.GetServiceCommands(anode)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Change Suricata Table"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "GetServiceCommands")    
    return data,err
}