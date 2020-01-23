package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "owlhnode/validation"
)

type ChangecontrolController struct {
    beego.Controller
}

// @Title GetChangeControlNode
// @Description Get changeControl database values
// @Success 200 {object} models.changecontrol
// @router / [get]
func (n *ChangecontrolController) GetChangeControlNode() {  
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{    
        data, err := models.GetChangeControlNode()
        n.Data["json"] = data
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }  
    n.ServeJSON()
}