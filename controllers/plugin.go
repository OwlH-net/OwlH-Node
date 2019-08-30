package controllers

import (
	"github.com/astaxie/beego"
	"owlhnode/models"
	"encoding/json"
)

type PluginController struct {
	beego.Controller
}

// @Title ChangeServiceStatus
// @Description Change a specific plugin service status
// @router /ChangeServiceStatus [put]
func (n *PluginController) ChangeServiceStatus() {
	var anode map[string]string
	json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
	err := models.ChangeServiceStatus(anode)
	n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title ChangeMainServiceStatus
// @Description Change a specific plugin service status
// @router /ChangeMainServiceStatus [put]
func (n *PluginController) ChangeMainServiceStatus() {
	var anode map[string]string
	json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
	err := models.ChangeMainServiceStatus(anode)
	n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}