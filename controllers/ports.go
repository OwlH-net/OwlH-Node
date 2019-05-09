package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"owlhnode/models"
)

type PortsController struct {
	beego.Controller
}

// @Title ShowPorts
// @Description get ports status
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
    logs.Info ("ports controller -> GET")
	err := models.ChangeMode()
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
    logs.Info ("ports controller -> GET")
	err := models.ChangeStatus()
	m.Data["json"] = map[string]string{"ack": "true"}
	if err != nil {
        logs.Info("ChangeStatus OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    m.ServeJSON()
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