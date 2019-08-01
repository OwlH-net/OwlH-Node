package controllers

import (
	"github.com/astaxie/beego"
	"owlhnode/models"
)

type PingController struct {
	beego.Controller
}

// @Title PingNode
// @Description get ping for node
// @router / [get]
func (n *PingController) PingNode() {
	n.Data["json"] = map[string]string{"ack": "true"}
    n.ServeJSON()
}

// @Title PingService
// @Description get ping for node
// @router /services [get]
func (n *PingController) PingService() {
	err := models.PingService()
	n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title DeployService
// @Description get ping for node
// @router /deployservice [put]
func (n *PingController) DeployService() {
	err := models.DeployService()
	n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}