package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"owlhnode/models"
	"owlhnode/validation"
	"encoding/json"
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
	token,err := models.CreateMasterToken(anode)
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
	permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), "none", n.Ctx.Input.Header("uuid"), "put")
	logs.Notice(permissions)
    if err != nil {
        logs.Error("AddUserFromMaster Error validating token from master")
        logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{ 
		var anode map[string]map[string]string
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		err := models.AddUserFromMaster(anode)
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
	permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), "none", n.Ctx.Input.Header("uuid"), "put")
	logs.Notice(permissions)
    if err != nil {
        logs.Error("AddRolesFromMaster Error validating token from master")
        logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
        logs.Error("Not enough permissions")
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{ 
		var anode map[string]map[string]string
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		err := models.AddRolesFromMaster(anode)
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
	permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), "none", n.Ctx.Input.Header("uuid"), "put")
	logs.Notice(permissions)
    if err != nil {
        logs.Error("AddGroupFromMaster Error validating token from master")
        logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{ 
		var anode map[string]map[string]string
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		err := models.AddGroupFromMaster(anode)
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
	permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), "none", n.Ctx.Input.Header("uuid"), "put")
	logs.Notice(permissions)
    if err != nil {
        logs.Error("AddUserGroupRolesFromMaster Error validating token from master")
        logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{ 
		var anode map[string]map[string]string
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		err := models.AddUserGroupRolesFromMaster(anode)
		n.Data["json"] = map[string]string{"ack": "true"}
		if err != nil {
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
		n.ServeJSON()
	}
}