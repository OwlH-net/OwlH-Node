package controllers

import (
    "github.com/astaxie/beego"
    "owlhnode/models"
    "encoding/json"
    "owlhnode/validation"
    "github.com/astaxie/beego/logs"
)

type PingController struct {
    beego.Controller
}

// @Title PingNode
// @Description get ping for node
// @router / [get]
func (n *PingController) PingNode() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{    
        n.Data["json"] = map[string]string{"ack": "true"}
    }
    n.ServeJSON()
}

// @Title UpdateNodeData
// @Description update node data
// @router /updateNode [put]
func (n *PingController) UpdateNodeData() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        logs.Info("ACTION -> PUT")
        logs.Info("CONTROLLER -> PING")
        logs.Info("ROUTER -> @router /updateNode [put]")
        for key := range anode {
            logs.Info("key -> "+key)
        }
        err := models.UpdateNodeData(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title PingService
// @Description get ping for node
// @router /services [get]
func (n *PingController) PingService() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        err := models.PingService()
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title DeployService
// @Description get ping for node
// @router /deployservice [put]
func (n *PingController) DeployService() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        logs.Info("ACTION -> PUT")
        logs.Info("CONTROLLER -> PING")
        logs.Info("ROUTER -> @router /deployservice [put]")
        err := models.DeployService(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title GetMainconfData
// @Description get ping for node
// @router /mainconf [get]
func (n *PingController) GetMainconfData() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        data,err := models.GetMainconfData()
        n.Data["json"] = data
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title PingPluginsNode
// @Description PingPluginsNode status
// @Success 200 {object} models.ports
// @router /PingPluginsNode [get]
func (n *PingController) PingPluginsNode() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        data, err := models.PingPluginsNode()
        n.Data["json"] = data
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}