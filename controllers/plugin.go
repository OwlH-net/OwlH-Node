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
    anode["action"] = "PUT"
    anode["controller"] = "PLUGIN"
    anode["router"] = "@router /ChangeServiceStatus [put]"
    err := models.ChangeServiceStatus(anode)
    // n.Data["json"] = err
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
    anode["action"] = "PUT"
    anode["controller"] = "PLUGIN"
    anode["router"] = "@router /ChangeMainServiceStatus [put]"
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
    anode["action"] = "DELETE"
    anode["controller"] = "PLUGIN"
    anode["router"] = "@router /deleteService [delete]"
    err := models.DeleteService(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title AddPluginService()
// @Description Add new service like Suricata or Zeek
// @Success 200 {object} models.suricata
// @router /addService [put]
func (n *PluginController) AddPluginService() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "PLUGIN"
    anode["router"] = "@router /addService [put]"
    err := models.AddPluginService(anode)

    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title SaveSuricataInterface
// @Description Change a specific plugin service status
// @router /SaveSuricataInterface [put]
func (n *PluginController) SaveSuricataInterface() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "PLUGIN"
    anode["router"] = "@router /SaveSuricataInterface [put]"
    err := models.SaveSuricataInterface(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title DeployStapService
// @Description Change a specific plugin service status
// @router /deployStapService [put]
func (n *PluginController) DeployStapService() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "PLUGIN"
    anode["router"] = "@router /deployStapService [put]"
    err := models.DeployStapService(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title StopStapService
// @Description Change a specific plugin service status
// @router /stopStapService [put]
func (n *PluginController) StopStapService() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "PLUGIN"
    anode["router"] = "@router /stopStapService [put]"
    err := models.StopStapService(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title ModifyStapValues
// @Description Change a specific plugin service status
// @router /modifyStapValues [put]
func (n *PluginController) ModifyStapValues() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "PLUGIN"
    anode["router"] = "@router /modifyStapValues [put]"
    err := models.ModifyStapValues(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title ChangeSuricataTable
// @Description Change a specific plugin service status
// @router /changeSuricataTable [put]
func (n *PluginController) ChangeSuricataTable() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "PLUGIN"
    anode["router"] = "@router /changeSuricataTable [put]"
    err := models.ChangeSuricataTable(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}