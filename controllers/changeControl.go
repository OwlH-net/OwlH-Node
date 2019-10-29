package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
)

type ChangecontrolController struct {
    beego.Controller
}

// @Title GetChangeControlNode
// @Description Get changeControl database values
// @Success 200 {object} models.changecontrol
// @router / [get]
func (m *ChangecontrolController) GetChangeControlNode() {    
    data, err := models.GetChangeControlNode()
    m.Data["json"] = data
    if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}