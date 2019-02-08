package controllers

import (
	"owlhnode/models"
	//"encoding/json"
	"strconv"
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
	mstatus := models.GetZeek()
	m.Data["json"] = map[string]string{"status": strconv.FormatBool(mstatus)}
    m.ServeJSON()
}

