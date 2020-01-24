package controllers

import (
    "owlhnode/models"
    "encoding/json"
    "github.com/astaxie/beego"
    "owlhnode/validation"
    "github.com/astaxie/beego/logs"
)

type GroupController struct {
    beego.Controller
}

// @Title SyncSuricataGroupValues
// @Description get Suricata group values
// @Success 200 {object} models.suricata
// @router /sync [put]
func (n *GroupController) SyncSuricataGroupValues() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Group Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.SyncSuricataGroupValues(anode)
        
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Error("SyncSuricataGroupValues controller -> GET -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title SuricataGroupService
// @Description Suricata start/stop for group node
// @Success 200 {object} models.suricata
// @router /suricata [put]
func (n *GroupController) SuricataGroupService() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Group Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.SuricataGroupService(anode)
    
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Error("SuricataGroupService controller -> GET -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}