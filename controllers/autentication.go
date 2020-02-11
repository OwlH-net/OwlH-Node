package controllers

import (
	"github.com/astaxie/beego"
	"owlhnode/models"
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