package controllers

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "owlhnode/models"
    "owlhnode/validation"
    "owlhnode/knownports"
    "encoding/json"
)

type PortsController struct {
    beego.Controller
}

// @Title PingPorts
// @Description PingPorts status
// @Success 200 {object} models.ports
// @router /PingPorts [get]
func (n *PortsController) PingPorts() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"PingPorts"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("uuid"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{       
        data, err := models.PingPorts()
        n.Data["json"] = data
        if err != nil {
            logs.Info("PingPorts OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title ShowPorts
// @Description get ports
// @Success 200 {object} models.ports
// @router / [get]
func (n *PortsController) ShowPorts() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"ShowPorts"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("uuid"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{      
        logs.Info ("ports controller -> GET")
        data,err := models.ShowPorts()
        n.Data["json"] = data
        if err != nil {
            logs.Info("ShowPorts OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title ChangeMode
// @Description put new mode
// @Success 200 {object} models.ports
// @router /mode [put]
func (n *PortsController) ChangeMode() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"ChangeMode"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("uuid"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{        
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "PORTS"
        anode["router"] = "@router /mode [put]"
    
        err := models.ChangeMode(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Info("ChangeMode OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title ChangeStatus
// @Description put new status
// @Success 200 {object} models.ports
// @router /status [put]
func (n *PortsController) ChangeStatus() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"ChangeStatus"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("uuid"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{      
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "PORTS"
        anode["router"] = "@router /status [put]"    
        anode["plugin"] = "knownports"
        logs.Info ("ports controller -> GET")
        err := models.ChangeStatus(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Info("ChangeStatus OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }else{
            knownports.Init()
        }
    }
    n.ServeJSON()    
    
}

// @Title DeletePorts
// @Description delete ports
// @Success 200 {object} models.ports
// @router /delete [put]
func (n *PortsController) DeletePorts() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"DeletePorts"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("uuid"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{          
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "PORTS"
        anode["router"] = "@router /delete [put]"
        err := models.DeletePorts(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Info("DeletePorts OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}


// @Title DeleteAllPorts
// @Description delete all ports
// @Success 200 {object} models.ports
// @router /deleteAll [put]
func (n *PortsController) DeleteAllPorts() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"DeleteAllPorts"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("uuid"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{      
        var anode map[string]string
        anode["action"] = "PUT"
        anode["controller"] = "PORTS"
        anode["router"] = "@router /mode [put]"
        err := models.DeleteAllPorts(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Info("DeletePorts OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}