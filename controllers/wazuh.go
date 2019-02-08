package controllers

import (
	"owlhnode/models"
	//"encoding/json"
	"strconv"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
)

type WazuhController struct {
	beego.Controller
}

// @Title GetWazuh
// @Description get Wazuh status
// @Success 200 {object} models.wazuh
// @router / [get]
func (m *WazuhController) Get() {
    logs.Info ("Wazuh controller -> GET")
	mstatus := models.GetWazuh()
	m.Data["json"] = map[string]string{"status": strconv.FormatBool(mstatus)}
    m.ServeJSON()
}

