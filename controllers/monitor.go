package controllers

import (
	"owlhnode/models"
	"owlhnode/monitor"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type MonitorController struct {
	beego.Controller
}


// @Title GetLastStats
// @Description get last node stats 
// @Success 200 {object} models.monitor
// @router / [get]
func (m *AnalyzerController) GetLastStatus() {	
    data := models.GetNodeStats()
    m.Data["json"] = data
    m.ServeJSON()
}
