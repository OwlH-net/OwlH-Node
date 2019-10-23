package models 

import (
    "owlhnode/deploy"
)

func DeployNode(anode map[string]string)(err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("DEPLOY - DeployNode")
    for key :=range anode {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    
    err = deploy.DeployNode(anode)
    return err
}

func CheckDeployFiles()(anode map[string]string) {
    anode = deploy.CheckDeployFiles()
    return anode
}