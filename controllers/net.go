package controllers

import (
	"github.com/astaxie/beego"
	"owlhnode/models"
)

type NetController struct {
	beego.Controller
}

// @Title PingNode
// @Description get ping for node
// @router / [get]
func (m *NetController) GetNetworkData() {
	data,err := models.GetNetworkData()
	m.Data["json"] = data
	if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    m.ServeJSON()
}