package controllers

import (
    "owlhnode/models"
    "encoding/json"
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

// @Title AddMonitorFile
// @Description Add file to monitor
// @Success 200 {object} models.monitor
// @router /addFile [post]
func (m *MonitorController) AddMonitorFile() {	
	var anode map[string]string
    json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
	anode["action"] = "POST"
    anode["controller"] = "MONITOR"
    anode["router"] = "@router /addFile [post]"

	err := models.AddMonitorFile(anode)
	m.Data["json"] = map[string]string{"ack": "true"}
	if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    
    m.ServeJSON()
}

// @Title PingMonitorFiles
// @Description get monitor file stats
// @Success 200 {object} models.monitor
// @router /pingMonitorFiles [get]
func (m *MonitorController) PingMonitorFiles() {	
    data,err := models.PingMonitorFiles()
    m.Data["json"] = data
    if err != nil {m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}}

    m.ServeJSON()
}

// @Title DeleteMonitorFile
// @Description Add file to monitor
// @Success 200 {object} models.monitor
// @router /deleteFile [delete]
func (m *MonitorController) DeleteMonitorFile() {	
	var anode map[string]string
    json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
	anode["action"] = "DELETE"
    anode["controller"] = "MONITOR"
    anode["router"] = "@router /deleteFile [delete]"
    
	err := models.DeleteMonitorFile(anode)
	m.Data["json"] = map[string]string{"ack": "true"}
	if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    
    m.ServeJSON()
}