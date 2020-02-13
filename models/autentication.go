package models

import (
    "owlhnode/autentication"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)

func CreateMasterToken(login map[string]string) (token string, err error) {
    logs.Info("============")
    logs.Info("AUTENTICATION - CreateMasterToken")
    cc := make(map[string]string)
    
    token,err = autentication.CreateMasterToken(login)
    
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Create master token"


    changecontrol.InsertChangeControl(cc)
    return token, err
}

func AddUserFromMaster(user map[string]map[string]string) (err error) {
    logs.Info("============")
    logs.Info("AUTENTICATION - AddUserFromMaster")
    cc := make(map[string]string)
    
    err = autentication.AddUserFromMaster(user)
    
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Add users from Master"


    changecontrol.InsertChangeControl(cc)
    return err
}