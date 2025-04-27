package controllers

import (
	"encoding/json"

	"github.com/OwlH-net/OwlH-Node/models"
	"github.com/OwlH-net/OwlH-Node/validation"
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
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"DeployNode"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		var anode map[string]string
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		anode["action"] = "PUT"
		anode["controller"] = "DEPLOY"
		anode["router"] = "@router / [put]"
		err := models.DeployNode(anode, n.Ctx.Input.Header("user"))
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
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"Deploy"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		anode := models.CheckDeployFiles(n.Ctx.Input.Header("user"))
		n.Data["json"] = anode
	}
	n.ServeJSON()
}
