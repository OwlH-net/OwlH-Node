package controllers

import (
    "owlhnode/models"
    "owlhnode/analyzer"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "encoding/json"
    "owlhnode/validation"
)

type AnalyzerController struct {
    beego.Controller
}

// @Title PingAnalyzer
// @Description PingAnalyzer status
// @Success 200 {object} models.analyzer
// @router /pingAnalyzer [get]
func (n *AnalyzerController) PingAnalyzer() {  
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("PingAnalyzer Error validating token from master")
logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{    
        data, err := models.PingAnalyzer()
        n.Data["json"] = data
        if err != nil {
            logs.Error("PingAnalyzer OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }  
    n.ServeJSON()
}

// @Title ChangeAnalyzerStatus
// @Description ChangeAnalyzerStatus status
// @Success 200 {object} models.analyzer
// @router /changeAnalyzerStatus [put]
func (n *AnalyzerController) ChangeAnalyzerStatus() { 
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("ChangeAnalyzerStatus Error validating token from master")
logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{    
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "ANALYZER"
        anode["router"] = "@router /changeAnalyzerStatus [put]"
        err := models.ChangeAnalyzerStatus(anode)
    
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Error("ChangeAnalyzerStatus OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }else{
            analyzer.Init()
        }
    }   
    
    n.ServeJSON()
}

// @Title SyncAnalyzer
// @Description SyncAnalyzer status
// @Success 200 {object} models.analyzer
// @router /sync [put]
func (n *AnalyzerController) SyncAnalyzer() { 
    
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("SyncAnalyzer Error validating token from master")
logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{    
        var anode map[string][]byte
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.SyncAnalyzer(anode)
    
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Error("SyncAnalyzer OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }   
    n.ServeJSON()
}