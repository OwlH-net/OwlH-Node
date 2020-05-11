package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "owlhnode/validation"
)

type CollectorController struct {
    beego.Controller
}

// @Title PlayCollector
// @Description Play collector
// @Success 200 {object} models.Collector
// @Failure 403 body is empty
// @router /play [get]
func (n *CollectorController) PlayCollector() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"PlayCollector"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{    
        err := models.PlayCollector()
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Error("COLLECTOR CREATE -> error: %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title StopCollector
// @Description Stop collector
// @Success 200 {object} models.Collector
// @Failure 403 body is empty
// @router /stop [get]
func (n *CollectorController) StopCollector() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"StopCollector"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{    
        err := models.StopCollector()
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Error("COLLECTOR CREATE -> error: %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title ShowCollector
// @Description Show collector
// @Success 200 {object} models.Collector
// @Failure 403 body is empty
// @router /show [get]
func (n *CollectorController) ShowCollector() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"ShowCollector"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{    
        data, err := models.ShowCollector()
        n.Data["json"] =  map[string]string{"data": data}
    
        if err != nil {
            logs.Error("COLLECTOR CREATE -> error: %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}