package models 

import (
    "owlhnode/deploy"
    "owlhnode/changeControl"
    "errors"
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
    //changecontrol.ChangeControlInsertData(err, "DeployNode")    
    return err
}

func CheckDeployFiles()(anode map[string]string) {
    anode = deploy.CheckDeployFiles()
    var err error
    if anode == nil { err = errors.New("No files") } else { err = nil }
    //changecontrol.ChangeControlInsertData(err, "CheckDeployFiles")    
    return anode
}