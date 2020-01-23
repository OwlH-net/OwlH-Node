package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "owlhnode/validation"
    "encoding/json"
)

type IncidentslController struct {
    beego.Controller
}

// @Title GetIncidentsNode
// @Description Get incidents database values
// @Success 200 {object} models.incidents
// @router / [get]
func (n *IncidentslController) GetIncidentsNode() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        data, err := models.GetIncidentsNode()
        n.Data["json"] = data
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }    
    n.ServeJSON()
}

// @Title PutIncidentNode
// @Description Put incidents into database
// @Success 200 {object} models.monitor
// @router / [post]
func (n *MonitorController) PutIncidentNode() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "POST"
        anode["controller"] = "INCIDENTS"
        anode["router"] = "@router / [post]"
        err := models.PutIncidentNode(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }    
    
    n.ServeJSON()
}