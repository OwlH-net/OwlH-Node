package controllers

import (
	"encoding/json"

	"github.com/OwlH-net/OwlH-Node/models"
	"github.com/OwlH-net/OwlH-Node/validation"
	"github.com/astaxie/beego"
)

type IncidentslController struct {
	beego.Controller
}

// @Title GetIncidentsNode
// @Description Get incidents database values
// @Success 200 {object} models.incidents
// @router / [get]
func (n *IncidentslController) GetIncidentsNode() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"GetIncidentsNode"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		data, err := models.GetIncidentsNode(n.Ctx.Input.Header("user"))
		n.Data["json"] = data
		if err != nil {
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
	}
	n.ServeJSON()
}

// @Title PutIncidentNode
// @Description Put incidents into database
// @Success 200 {object} models.monitor
// @router / [post]
func (n *MonitorController) PutIncidentNode() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"PutIncidentNode"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		var anode map[string]string
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		anode["action"] = "POST"
		anode["controller"] = "INCIDENTS"
		anode["router"] = "@router / [post]"
		err := models.PutIncidentNode(anode, n.Ctx.Input.Header("user"))
		n.Data["json"] = map[string]string{"ack": "true"}
		if err != nil {
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
	}

	n.ServeJSON()
}
