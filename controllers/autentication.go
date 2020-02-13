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
		n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    n.ServeJSON()
}

// @Title AddUserFromMaster
// @Description Get Master user
// @router /addUser [put]
func (n *AutenticationController) AddUserFromMaster() {
	err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("user"))
    if err != nil {
        logs.Error("AddUserFromMaster Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
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