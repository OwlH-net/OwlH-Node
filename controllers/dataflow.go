package controllers

import (
    "owlhnode/models"
    "encoding/json"
    "owlhnode/validation"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
)

type DataflowController struct {
    beego.Controller
}

// @Title ChangeDataflowValues
// @Description Make a deploy for selected action
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /changeDataflowValues [put]
func (n *DataflowController) ChangeDataflowValues() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)    
        err := models.ChangeDataflowValues(anode)
        anode["action"] = "PUT"
        anode["controller"] = "DATAFLOW"
        anode["router"] = "@router /changeDataflowValues [put]"
    
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title LoadDataflowValues
// @Description Load data flow values for
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /loadDataflowValues [get]
func (n *DataflowController) LoadDataflowValues() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        data,err := models.LoadDataflowValues()
        n.Data["json"] = data
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title SaveSocketToNetwork
// @Description Save socket information to Network at node dataflow
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /saveSocketToNetwork [put]
func (n *DataflowController) SaveSocketToNetwork() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)    
        anode["action"] = "PUT"
        anode["controller"] = "DATAFLOW"
        anode["router"] = "@router /saveSocketToNetwork [put]"    
        err := models.SaveSocketToNetwork(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title SaveNewLocal
// @Description Save socket information to Network at node dataflow
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /saveNewLocal [put]
func (n *DataflowController) SaveNewLocal() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)    
        anode["action"] = "PUT"
        anode["controller"] = "DATAFLOW"
        anode["router"] = "@router /saveNewLocal [put]"    
        err := models.SaveNewLocal(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title SaveVxLAN
// @Description Save socket information to Network at node dataflow
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /saveVxLAN [put]
func (n *DataflowController) SaveVxLAN() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)    
        anode["action"] = "PUT"
        anode["controller"] = "DATAFLOW"
        anode["router"] = "@router /saveVxLAN [put]"    
        err := models.SaveVxLAN(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title SaveSocketToNetworkSelected
// @Description Save socket information to Network at node dataflow
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /saveSocketToNetworkSelected [put]
func (n *DataflowController) SaveSocketToNetworkSelected() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)    
        anode["action"] = "PUT"
        anode["controller"] = "DATAFLOW"
        anode["router"] = "@router /saveSocketToNetworkSelected [put]"    
        err := models.SaveSocketToNetworkSelected(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title DeleteDataFlowValueSelected
// @Description Save socket information to Network at node dataflow
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /deleteDataFlowValueSelected [delete]
func (n *DataflowController) DeleteDataFlowValueSelected() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "delete")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)    
        anode["action"] = "DELETE"
        anode["controller"] = "DATAFLOW"
        anode["router"] = "@router /deleteDataFlowValueSelected [delete]"    
        err := models.DeleteDataFlowValueSelected(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}