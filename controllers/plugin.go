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

// @Title DeleteService
// @Description delete a specific plugin service
// @router /deleteService [delete]
func (n *PluginController) DeleteService() {
	var anode map[string]string
	json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
	err := models.DeleteService(anode)
	n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title AddPluginService()
// @Description Add new Suricata service
// @Success 200 {object} models.suricata
// @router /addService [put]
func (n *PluginController) AddPluginService() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    err := models.AddPluginService(anode)

    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}