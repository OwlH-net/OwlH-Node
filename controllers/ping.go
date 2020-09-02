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
    // permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    // if err != nil {
    //     logs.Error("Error validating token from master")
    //     n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "nodeToken":"none"}
    // }else{    
        n.Data["json"] = map[string]string{"ack": "true"}
    // }
    n.ServeJSON()
}

// @Title UpdateNodeData
// @Description update node data
// @router /updateNode [put]
func (n *PingController) UpdateNodeData() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"UpdateNodeData"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        logs.Info("ACTION -> PUT")
        logs.Info("CONTROLLER -> PING")
        logs.Info("ROUTER -> @router /updateNode [put]")
        for key := range anode {
            logs.Info("key -> "+key)
        }
        err := models.UpdateNodeData(anode, n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"PingService"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        err := models.PingService(n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"DeployService"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        logs.Info("ACTION -> PUT")
        logs.Info("CONTROLLER -> PING")
        logs.Info("ROUTER -> @router /deployservice [put]")
        err := models.DeployService(anode, n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"GetMainconfData"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        data,err := models.GetMainconfData(n.Ctx.Input.Header("user"))
        n.Data["json"] = data
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title PingPluginsNode
// @Description PingPluginsNode status
// @Success 200 {object} models.ping
// @router /PingPluginsNode [get]
func (n *PingController) PingPluginsNode() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        var errorResponse = map[string]map[string]string{}
        errorResponse["hasError"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.Data["json"] = errorResponse
        n.ServeJSON()
        return
    }    
    permissions := []string{"PingPluginsNode"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        var errorResponse = map[string]map[string]string{}
        errorResponse["hasError"] = map[string]string{"ack": "false","permissions":"none", "error": "Not enough permissions"}
        n.Data["json"] = errorResponse
    }else{         
        data, err := models.PingPluginsNode(n.Ctx.Input.Header("user"))
        n.Data["json"] = data
        if err != nil {
            var errorResponse = map[string]map[string]string{}
            errorResponse["hasError"] = map[string]string{"ack": "false", "error": err.Error()}
            n.Data["json"] = errorResponse
        }
    }
    n.ServeJSON()
}

// @Title SaveNodeInformation
// @Description Save node information
// @Success 200 {object} models.ping
// @router /saveNodeInformation [put]
func (n *PingController) SaveNodeInformation() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"SaveNodeInformation"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        logs.Info("ACTION -> PUT")
        logs.Info("CONTROLLER -> PING")
        logs.Info("ROUTER -> @router /SaveNodeInformation [put]")
        err := models.SaveNodeInformation(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title DeleteNode
// @Description Save node information
// @Success 200 {object} models.ping
// @router / [delete]
func (n *PingController) DeleteNode() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"DeleteNode"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{              
        err := models.DeleteNode(n.Ctx.Input.Header("uuid"), n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}