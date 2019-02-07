package controllers

import (
	"owlhnode/models"
	//"encoding/json"
	"strconv"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
)

type SuricataController struct {
	beego.Controller
}

// @Title GetSuricata
// @Description get Surucata status
// @Success 200 {object} models.suricata
// @router / [get]
func (m *SuricataController) Get() {
    logs.Info ("Suricata controller -> GET")
	mstatus := models.GetSuricata()
	m.Data["json"] = map[string]string{"status": strconv.FormatBool(mstatus)}
    m.ServeJSON()
}

