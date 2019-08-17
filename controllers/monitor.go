package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
)

type MonitorController struct {
	beego.Controller
}


// @Title GetLastStats
// @Description get last node stats 
// @Success 200 {object} models.monitor
// @router / [get]
func (m *MonitorController) GetLastStatus() {	
    data := models.GetNodeStats()
    m.Data["json"] = data
    m.ServeJSON()
}
