package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
    "encoding/json"
)

type IncidentslController struct {
    beego.Controller
}

// @Title GetIncidentsNode
// @Description Get incidents database values
// @Success 200 {object} models.incidents
// @router / [get]
func (m *IncidentslController) GetIncidentsNode() {    
    data, err := models.GetIncidentsNode()
    m.Data["json"] = data
    if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}

// @Title PutIncidentNode
// @Description Put incidents into database
// @Success 200 {object} models.monitor
// @router / [post]
func (m *MonitorController) PutIncidentNode() {    
    var anode map[string]string
    json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
    anode["action"] = "POST"
    anode["controller"] = "INCIDENTS"
    anode["router"] = "@router / [post]"
    err := models.PutIncidentNode(anode)
    m.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    
    m.ServeJSON()
}