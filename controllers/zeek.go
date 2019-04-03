package controllers

import (
	"owlhnode/models"
	//"encoding/json"
	// "strconv"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
)

type ZeekController struct {
	beego.Controller
}

// @Title GetZeek
// @Description get Zeek status
// @Success 200 {object} models.zeek
// @router / [get]
func (m *ZeekController) Get() {
    logs.Info ("Zeek controller -> GET")
	mstatus,err := models.GetZeek()
	m.Data["json"] = mstatus
	if err != nil {
        logs.Info("GetWazuh OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}

// @Title RunZeek
// @Description Run zeek system
// @Success 200 {object} models.zeek
// @Failure 403 body is empty
// @router /RunZeek [put]
func (n *ZeekController) RunZeek() {
    logs.Info("RunZeek -> In")
    data,err := models.RunZeek()
    n.Data["json"] = data
    if err != nil {
        logs.Info("RunZeek OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("RunZeek -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title StopZeek
// @Description Run zeek system
// @Success 200 {object} models.zeek
// @Failure 403 body is empty
// @router /StopZeek [put]
func (n *ZeekController) StopZeek() {
    logs.Info("StopZeek -> In")
    data,err := models.StopZeek()
    n.Data["json"] = data
    if err != nil {
        logs.Info("StopZeek OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("StopZeek -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}