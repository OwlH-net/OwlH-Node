package models 

import (
    "owlhnode/deploy"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func DeployNode(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("DEPLOY - DeployNode")
    for key :=range anode {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    
    err = deploy.DeployNode(anode)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Deploy node"

    changecontrol.InsertChangeControl(cc)

    //changecontrol.ChangeControlInsertData(err, "DeployNode")    
    return err
}

func CheckDeployFiles()(anode map[string]string) {
    anode = deploy.CheckDeployFiles()
    // if anode == nil { err = errors.New("No files") } else { err = nil }
    //changecontrol.ChangeControlInsertData(err, "CheckDeployFiles")    
    return anode
}