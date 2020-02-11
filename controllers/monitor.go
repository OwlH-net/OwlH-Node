package controllers

import (
    "owlhnode/models"
    "encoding/json"
    "owlhnode/validation"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
)

type MonitorController struct {
    beego.Controller
}


// @Title GetLastStats
// @Description get last node stats 
// @Success 200 {object} models.monitor
// @router / [get]
func (n *MonitorController) GetLastStatus() {  
    data := models.GetNodeStats()
    n.Data["json"] = data
    
    n.ServeJSON()
}

// @Title AddMonitorFile
// @Description Add file to monitor
// @Success 200 {object} models.monitor
// @router /addFile [post]
func (n *MonitorController) AddMonitorFile() { 
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("user"))
    if err != nil {
        logs.Error("AddMonitorFile Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "POST"
        anode["controller"] = "MONITOR"
        anode["router"] = "@router /addFile [post]"
    
        err := models.AddMonitorFile(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }   
    
    n.ServeJSON()
}

// @Title PingMonitorFiles
// @Description get monitor file stats
// @Success 200 {object} models.monitor
// @router /pingMonitorFiles [get]
func (n *MonitorController) PingMonitorFiles() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("user"))
    if err != nil {
        logs.Error("PingMonitorFiles Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        data,err := models.PingMonitorFiles()
        n.Data["json"] = data
        if err != nil {n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}}
    }    

    n.ServeJSON()
}

// @Title DeleteMonitorFile
// @Description Add file to monitor
// @Success 200 {object} models.monitor
// @router /deleteFile [delete]
func (n *MonitorController) DeleteMonitorFile() { 
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("user"))
    if err != nil {
        logs.Error("DeleteMonitorFile Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "DELETE"
        anode["controller"] = "MONITOR"
        anode["router"] = "@router /deleteFile [delete]"
        
        err := models.DeleteMonitorFile(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }   
    
    n.ServeJSON()
}