package models 

import (
    "owlhnode/deploy"
)

func DeployNode(anode map[string]string)(err error) {
    err = deploy.DeployNode(anode)
    return err
}

func CheckDeployFiles()(anode map[string]string) {
    anode = deploy.CheckDeployFiles()
    return anode
}