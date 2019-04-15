package controllers

import (
	"github.com/astaxie/beego"
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