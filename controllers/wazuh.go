package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "encoding/json"
    "owlhnode/validation"
)

type WazuhController struct {
    beego.Controller
}

// @Title GetWazuh
// @Description get Wazuh status
// @Success 200 {object} models.wazuh
// @router / [get]
func (n *WazuhController) Get() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"GetWazuh"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        logs.Info ("Wazuh controller -> GET")
        mstatus, err := models.GetWazuh()
        n.Data["json"] = mstatus
        if err != nil {
            logs.Info("GetWazuh OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title RunWazuh
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /RunWazuh [put]
func (n *WazuhController) RunWazuh() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"RunWazuh"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{          
        logs.Info("RunWazuh -> In")
        var anode map[string]string
        anode["action"] = "PUT"
        anode["controller"] = "WAZUH"
        anode["router"] = "@router /RunWazuh [put]"
        logs.Info("============")
        logs.Info("WAZUH - RunWazuh")
        for key :=range anode {
            logs.Info(key +" -> "+anode[key])
        }
        data,err := models.RunWazuh()
        n.Data["json"] = data
        if err != nil {
            logs.Info("RunWazuh OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        logs.Info("RunWazuh -> OUT -> %s", n.Data["json"])
    }
    n.ServeJSON()
}

// @Title StopWazuh
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /StopWazuh [put]
func (n *WazuhController) StopWazuh() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"StopWazuh"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{        
        logs.Info("StopWazuh -> In")
        var anode map[string]string
        anode["action"] = "PUT"
        anode["controller"] = "WAZUH"
        anode["router"] = "@router /StopWazuh [put]"
        logs.Info("============")
        logs.Info("WAZUH - StopWazuh")
        for key :=range anode {
            logs.Info(key +" -> "+anode[key])
        }
        data,err := models.StopWazuh()
        n.Data["json"] = data
        if err != nil {
            logs.Info("StopWazuh OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        logs.Info("StopWazuh -> OUT -> %s", n.Data["json"])
    }
    n.ServeJSON()
}

// @Title PingWazuhFiles
// @Description get Wazuh status
// @Success 200 {object} models.wazuh
// @router /pingWazuhFiles [get]
func (n *WazuhController) PingWazuhFiles() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        var errorResponse = map[string]map[string]string{}
        errorResponse["hasError"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.Data["json"] = errorResponse
        n.ServeJSON()
        return
    }    
    permissions := []string{"PingWazuhFiles"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        var errorResponse = map[string]map[string]string{}
        errorResponse["hasError"] = map[string]string{"ack": "false","permissions":"none", "error": "Not enough permissions"}
        n.Data["json"] = errorResponse
    }else{       
        files, err := models.PingWazuhFiles()
        n.Data["json"] = files
        if err != nil {
            var errorResponse = map[string]map[string]string{}
            errorResponse["hasError"] = map[string]string{"ack": "false", "error": err.Error()}
            n.Data["json"] = errorResponse
        }
    }
    n.ServeJSON()
}

// @Title DeleteWazuhFile
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /deleteWazuhFile [delete]
func (n *WazuhController) DeleteWazuhFile() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"DeleteWazuhFile"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{      
        var anode map[string]interface{}
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "SURICATA"
        anode["router"] = "@router /StopSuricata [put]"
        err := models.DeleteWazuhFile(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        
        if err != nil {n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}}
    }
    n.ServeJSON()
}

// @Title AddWazuhFile
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /addWazuhFile [put]
func (n *WazuhController) AddWazuhFile() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"AddWazuhFile"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        var anode map[string]interface{}
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "SURICATA"
        anode["router"] = "@router /StopSuricata [put]"
        err := models.AddWazuhFile(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        
        if err != nil {n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}}
    }
    n.ServeJSON()
}

// @Title LoadFileLastLines
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /loadFileLastLines [put]
func (n *WazuhController) LoadFileLastLines() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"LoadFileLastLines"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{        
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "SURICATA"
        anode["router"] = "@router /StopSuricata [put]"
        data,err := models.LoadFileLastLines(anode)
        n.Data["json"] = data
        
        if err != nil {n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}}
    }
    n.ServeJSON()
}

// @Title SaveFileContentWazuh
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /saveFileContentWazuh [put]
func (n *WazuhController) SaveFileContentWazuh() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"SaveFileContentWazuh"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "SURICATA"
        anode["router"] = "@router /StopSuricata [put]"
    
        err := models.SaveFileContentWazuh(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        
        if err != nil {n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}}
    }
    n.ServeJSON()
}