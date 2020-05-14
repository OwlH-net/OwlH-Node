package controllers

import (
    "encoding/json"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "owlhnode/models"
    "owlhnode/validation"
)

type AutenticationController struct {
    beego.Controller
}

// @Title CreateMasterToken
// @Description Get Master Token
// @router / [put]
func (n *AutenticationController) CreateMasterToken() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    token, err := models.CreateMasterToken(anode, n.Ctx.Input.Header("user"))
    n.Data["json"] = map[string]string{"ack": "true", "token": token}
    if err != nil {
        logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title AddUserFromMaster
// @Description Add users from Master
// @router /addUser [put]
func (n *AutenticationController) AddUserFromMaster() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
        // }
        // permissions := []string{"AddUserFromMaster"}
        // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
        // if permissionsErr != nil || hasPermission == false {
        //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    } else {
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.AddUserFromMaster(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        n.ServeJSON()
    }
}

// @Title AddRolesFromMaster
// @Description Add roles from Master
// @router /addRole [put]
func (n *AutenticationController) AddRolesFromMaster() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
        // }
        // permissions := []string{"AddRolesFromMaster"}
        // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
        // if permissionsErr != nil || hasPermission == false {
        //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    } else {
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.AddRolesFromMaster(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        n.ServeJSON()
    }
}

// @Title AddGroupFromMaster
// @Description Add groups from Master
// @router /addGroup [put]
func (n *AutenticationController) AddGroupFromMaster() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
        // }
        // permissions := []string{"AddGroupFromMaster"}
        // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
        // if permissionsErr != nil || hasPermission == false {
        //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    } else {
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.AddGroupFromMaster(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        n.ServeJSON()
    }
}

// @Title AddUserGroupRolesFromMaster
// @Description Add userGroupRoles from Master
// @router /addUgr [put]
func (n *AutenticationController) AddUserGroupRolesFromMaster() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
        // }
        // permissions := []string{"AddUserGroupRolesFromMaster"}
        // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
        // if permissionsErr != nil || hasPermission == false {
        //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    } else {
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.AddUserGroupRolesFromMaster(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        n.ServeJSON()
    }
}

// @Title SyncRolePermissions
// @Description Add role permissions relation table from Master
// @router /addRolePerm [put]
func (n *AutenticationController) SyncRolePermissions() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
        // }
        // permissions := []string{"SyncRolePermissions"}
        // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
        // if permissionsErr != nil || hasPermission == false {
        //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    } else {
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.SyncRolePermissions(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        n.ServeJSON()
    }
}

// @Title SyncPermissions
// @Description Add permissions from Master
// @router /addPerm [put]
func (n *AutenticationController) SyncPermissions() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
        // }
        // permissions := []string{"SyncPermissions"}
        // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
        // if permissionsErr != nil || hasPermission == false {
        //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    } else {
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.SyncPermissions(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        n.ServeJSON()
    }
}

// @Title SyncRoleGroups
// @Description Add permissions from Master
// @router /addRoleGroups [put]
func (n *AutenticationController) SyncRoleGroups() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
        // }
        // permissions := []string{"SyncRoleGroups"}
        // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
        // if permissionsErr != nil || hasPermission == false {
        //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    } else {
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.SyncRoleGroups(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        n.ServeJSON()
    }
}
