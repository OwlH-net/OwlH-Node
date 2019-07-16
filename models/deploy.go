package models 

import (
    "owlhnode/deploy"
)

func Deploy(anode map[string]string)(err error) {
    err = deploy.Deploy(anode)
    return err
}

func CheckDeployFiles()(anode map[string]string) {
    anode = deploy.CheckDeployFiles()
    return anode
}