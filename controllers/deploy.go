package controllers

import (
    "owlhnode/models"
    "encoding/json"
    "owlhnode/validation"
    "github.com/astaxie/beego/logs"
    "github.com/astaxie/beego"
)

type DeployController struct {
    beego.Controller
}

// @Title DeployNode
// @Description Make a deploy for selected action
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router / [put]
func (n *DeployController) DeployNode() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Deploy Error validating token from master")
logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "DEPLOY"
        anode["router"] = "@router / [put]"
        err := models.DeployNode(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title Deploy
// @Description Get all the deploy file status
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router / [get]
func (n *DeployController) CheckDeployFiles() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("Deploy Error validating token from master")
logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        anode := models.CheckDeployFiles()
        n.Data["json"] = anode
    }
    n.ServeJSON()
}