package models

import (
    "owlhnode/autentication"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)

func CreateMasterToken() (token string, err error) {
    logs.Info("============")
    logs.Info("AUTENTICATION - CreateMasterToken")
    cc := make(map[string]string)
    
    token,err = autentication.CreateMasterToken()
    
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