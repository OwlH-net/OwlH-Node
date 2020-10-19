package models

import (
    "owlhnode/autentication"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)

func CreateMasterToken(login map[string]string, username string) (token string, err error) {
    logs.Info("============")
    logs.Info("AUTENTICATION - CreateMasterToken")
    cc := make(map[string]string)
    
    token,err = autentication.CreateMasterToken(login)
    if err!=nil { 
        logs.Error("NO TOKEN NO TOKEN NO TOKEN NO TOKEN NO TOKEN NO TOKEN NO TOKEN NO TOKEN ")
        logs.Error("NO TOKEN NO TOKEN NO TOKEN NO TOKEN NO TOKEN NO TOKEN NO TOKEN NO TOKEN ")
        logs.Error(err.Error())
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["username"] = username
    cc["actionDescription"] = "Create master token"

    changecontrol.InsertChangeControl(cc)
    return token, err
}

func AddUserFromMaster(user map[string]map[string]string, username string) (err error) {
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
    cc["username"] = username
    cc["actionDescription"] = "Add users from Master"

    changecontrol.InsertChangeControl(cc)
    return err
}

func AddRolesFromMaster(roles map[string]map[string]string, username string) (err error) {
    logs.Info("============")
    logs.Info("AUTENTICATION - AddRolesFromMaster")
    cc := make(map[string]string)
    
    err = autentication.AddRolesFromMaster(roles)
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["username"] = username
    cc["actionDescription"] = "Add roles from Master"

    changecontrol.InsertChangeControl(cc)
    return err
}

func AddGroupFromMaster(groups map[string]map[string]string, username string) (err error) {
    logs.Info("============")
    logs.Info("AUTENTICATION - AddGroupFromMaster")
    cc := make(map[string]string)
    
    err = autentication.AddGroupFromMaster(groups)
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["username"] = username
    cc["actionDescription"] = "Add groups from Master"

    changecontrol.InsertChangeControl(cc)
    return err
}

func AddUserGroupRolesFromMaster(ugr map[string]map[string]string, username string) (err error) {
    logs.Info("============")
    logs.Info("AUTENTICATION - AddUserGroupRolesFromMaster")
    cc := make(map[string]string)
    
    err = autentication.AddUserGroupRolesFromMaster(ugr)
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["username"] = username
    cc["actionDescription"] = "Add userGroupRoles from Master"

    changecontrol.InsertChangeControl(cc)
    return err
}

func SyncRolePermissions(rolePerm map[string]map[string]string, username string) (err error) {
    logs.Info("============")
    logs.Info("AUTENTICATION - SyncRolePermissions")
    cc := make(map[string]string)
    
    err = autentication.SyncRolePermissions(rolePerm)
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["username"] = username
    cc["actionDescription"] = "Add role permissions relation table from Master"

    changecontrol.InsertChangeControl(cc)
    return err
}

func SyncPermissions(perms map[string]map[string]string, username string) (err error) {
    logs.Info("============")
    logs.Info("AUTENTICATION - SyncPermissions")
    cc := make(map[string]string)
    
    err = autentication.SyncPermissions(perms)
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["username"] = username
    cc["actionDescription"] = "Add permissions from Master"

    changecontrol.InsertChangeControl(cc)
    return err
}

func SyncRoleGroups(perms map[string]map[string]string, username string) (err error) {
    logs.Info("============")
    logs.Info("AUTENTICATION - SyncRoleGroups")
    cc := make(map[string]string)
    
    err = autentication.SyncRoleGroups(perms)
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["username"] = username
    cc["actionDescription"] = "Add role groups from Master"

    changecontrol.InsertChangeControl(cc)
    return err
}