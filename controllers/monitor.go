package controllers

import (
    "owlhnode/models"
    "encoding/json"
    "owlhnode/validation"
    "github.com/astaxie/beego"
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"AddMonitorFile"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "POST"
        anode["controller"] = "MONITOR"
        anode["router"] = "@router /addFile [post]"
    
        err := models.AddMonitorFile(anode, n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"PingMonitorFiles"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        data,err := models.PingMonitorFiles(n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"DeleteMonitorFile"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "DELETE"
        anode["controller"] = "MONITOR"
        anode["router"] = "@router /deleteFile [delete]"
        
        err := models.DeleteMonitorFile(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }   
    
    n.ServeJSON()
}

// @Title ChangeRotationStatus
// @Description Change monitor rotation file status
// @Success 200 {object} models.monitor
// @router /changeRotationStatus [put]
func (n *MonitorController) ChangeRotationStatus() { 
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"ChangeRotationStatus"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "MONITOR"
        anode["router"] = "@router /changeRotationStatus [put]"
        
        err := models.ChangeRotationStatus(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }   
    
    n.ServeJSON()
}

// @Title EditRotation
// @Description Edit monitor rotation file parameters
// @Success 200 {object} models.monitor
// @router /editRotation [put]
func (n *MonitorController) EditRotation() { 
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"EditRotation"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "MONITOR"
        anode["router"] = "@router /editRotation [put]"
        
        err := models.EditRotation(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }   
    
    n.ServeJSON()
}