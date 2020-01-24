package controllers

import (
    "github.com/astaxie/beego"
    "owlhnode/models"
)

type AutenticationController struct {
    beego.Controller
}

// @Title CreateMasterToken
// @Description Get Master Token
// @router / [get]
func (n *AutenticationController) CreateMasterToken() {
	token,err := models.CreateMasterToken()
	n.Data["json"] = map[string]string{"ack": "true", "token": token}
	if err != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    n.ServeJSON()
}