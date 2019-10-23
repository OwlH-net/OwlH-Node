package models

import (
    "owlhnode/plugin"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)

func ChangeServiceStatus(anode map[string]string)(err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - ChangeServiceStatus")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

	err = plugin.ChangeServiceStatus(anode)
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
    var cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - ChangeMainServiceStatus")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
	err = plugin.ChangeMainServiceStatus(anode)
    return err
}

func DeleteService(anode map[string]string)(err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - DeleteService")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
	err = plugin.DeleteService(anode)
    return err
}

func AddPluginService(anode map[string]string) (err error) {
    var cc := anode
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
        anode["actionStatus"] = "error"
        anode["errorDescription"] = err.Error()
    }else{
        anode["actionStatus"] = "success"
    }
    anode["actionDescription"] = "Add plugin service"
    var controlError error
    controlError = changecontrol.InsertChangeControl(anode)
    if controlError!=nil { logs.Error("AddPluginService controlError: "+controlError.Error()) }
    if err != nil {return err}

    return err
}

func SaveSuricataInterface(anode map[string]string)(err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - SaveSuricataInterface")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
	err = plugin.SaveSuricataInterface(anode)
    return err
}

func DeployStapService(anode map[string]string)(err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - DeployStapService")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
	err = plugin.DeployStapService(anode)
    return err
}

func StopStapService(anode map[string]string)(err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - StopStapService")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
	err = plugin.StopStapService(anode)
    return err
}

func ModifyStapValues(anode map[string]string)(err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - ModifyStapValues")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
	err = plugin.ModifyStapValues(anode)
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
    var cc := anode
    logs.Info("============")
    logs.Info("PLUGIN - ChangeSuricataTable")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

	err = plugin.ChangeSuricataTable(anode)
    return err
}