package controllers

import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "owlhnode/models"
    "owlhnode/validation"
    "encoding/json"
)

type NetController struct {
    beego.Controller
}

// @Title GetNetworkData
// @Description get network data
// @router / [get]
func (n *NetController) GetNetworkData() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        values,err := models.GetNetworkData()
        
        n.Data["json"] = values
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title LoadNetworkValuesSelected
// @Description get network values selected by user
// @router /values [get]
func (n *NetController) LoadNetworkValuesSelected() {
    values := make(map[string]map[string]string)
    values["node"] = map[string]string{}
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        values["node"]["ack"] = "false"
        values["node"]["error"] = err.Error()
        values["node"]["token"] = "none"
        n.Data["json"] = values
    }else if !permissions{
        values["node"]["ack"] = "false"
        values["node"]["permissions"] = "none"
        n.Data["json"] = values
    }else{         
        values,err := models.LoadNetworkValuesSelected()
        
        n.Data["json"] = values
        if err != nil {
            values["node"]["ack"] = "false"
            values["node"]["error"] = err.Error()
            n.Data["json"] = values
        }
    }
    n.ServeJSON()
}

// @Title UpdateNetworkInterface
// @Description Make a deploy for selected action
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router / [put]
func (n *NetController) UpdateNetworkInterface() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "NET"
        anode["router"] = "@router / [put]"
        err := models.UpdateNetworkInterface(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// // @Title LoadNetworkValuesSuricata
// // @Description get network data
// // @router /loadNetworkValuesSuricata [get]
// func (n *NetController) LoadNetworkValuesSuricata() {
//     values,err := models.LoadNetworkValuesSuricata()
    
//     n.Data["json"] = values
//     if err != nil {
//         n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
//     }
//     n.ServeJSON()
// }