package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"owlhnode/models"
	"owlhnode/knownports"
	"encoding/json"
)

type PortsController struct {
	beego.Controller
}

// @Title PingPorts
// @Description PingPorts status
// @Success 200 {object} models.ports
// @router /PingPorts [get]
func (m *PortsController) PingPorts() {
	data, err := models.PingPorts()
	m.Data["json"] = data
	if err != nil {
        logs.Info("PingPorts OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    m.ServeJSON()
}

// @Title ShowPorts
// @Description get ports
// @Success 200 {object} models.ports
// @router / [get]
func (m *PortsController) ShowPorts() {
    logs.Info ("ports controller -> GET")
	data,err := models.ShowPorts()
	m.Data["json"] = data
	if err != nil {
        logs.Info("ShowPorts OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    m.ServeJSON()
}

// @Title ChangeMode
// @Description put new mode
// @Success 200 {object} models.ports
// @router /mode [put]
func (m *PortsController) ChangeMode() {
    var anode map[string]string
	json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "PORTS"
    anode["router"] = "@router /mode [put]"

	err := models.ChangeMode(anode)
	m.Data["json"] = map[string]string{"ack": "true"}
	if err != nil {
        logs.Info("ChangeMode OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    m.ServeJSON()
}

// @Title ChangeStatus
// @Description put new status
// @Success 200 {object} models.ports
// @router /status [put]
func (m *PortsController) ChangeStatus() {
	var anode map[string]string
	json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "PORTS"
    anode["router"] = "@router /status [put]"	
	anode["plugin"] = "knownports"
    logs.Info ("ports controller -> GET")
	err := models.ChangeStatus(anode)
	m.Data["json"] = map[string]string{"ack": "true"}
	if err != nil {
        logs.Info("ChangeStatus OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}else{
		knownports.Init()
	}
	m.ServeJSON()	
	
}

// @Title DeletePorts
// @Description delete ports
// @Success 200 {object} models.ports
// @router /delete [put]
func (m *PortsController) DeletePorts() {
	var anode map[string]string
	json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "PORTS"
    anode["router"] = "@router /delete [put]"
	err := models.DeletePorts(anode)
	m.Data["json"] = map[string]string{"ack": "true"}
	if err != nil {
        logs.Info("DeletePorts OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    m.ServeJSON()
}


// @Title DeleteAllPorts
// @Description delete all ports
// @Success 200 {object} models.ports
// @router /deleteAll [put]
func (m *PortsController) DeleteAllPorts() {
    anode["action"] = "PUT"
    anode["controller"] = "PORTS"
    anode["router"] = "@router /mode [put]"
	err := models.DeleteAllPorts(anode)
	m.Data["json"] = map[string]string{"ack": "true"}
	if err != nil {
        logs.Info("DeletePorts OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    m.ServeJSON()
}