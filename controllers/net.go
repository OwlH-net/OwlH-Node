package controllers

import (
	"github.com/astaxie/beego"
	"owlhnode/models"
	"encoding/json"
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

// @Title UpdateNetworkInterface
// @Description Make a deploy for selected action
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router / [put]
func (n *NetController) UpdateNetworkInterface() {
    var anode map[string]string
	json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    err := models.UpdateNetworkInterface(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}