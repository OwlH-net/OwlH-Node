package controllers

import (
	// "owlhnode/models"
	// "encoding/json"
	// "strconv"
	"github.com/astaxie/beego"
    // "github.com/astaxie/beego/logs"
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