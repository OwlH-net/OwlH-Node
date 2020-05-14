package controllers

import (
    "encoding/json"
    "github.com/astaxie/beego"
    // "github.com/astaxie/beego/logs"
    "owlhnode/models"
    "owlhnode/pcap"
    // "owlhnode/validation"
)

type HwaddmngController struct {
    beego.Controller
}

// @Title AddMacIp
// @Description Add MAC and IP to Knownmacs
// @router / [post]
func (n *HwaddmngController) AddMacIp() {
	// errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    // if errToken != nil {
    //     n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
    //     n.ServeJSON()
    //     return
    // }    
    // permissions := []string{"AddMacIp"}
    // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    // if permissionsErr != nil || hasPermission == false {
    //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    // }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "POST"
        anode["controller"] = "hwaddmng"
        anode["router"] = "@router / [post]"

        err := models.AddMacIp(anode)
    
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    // }
    n.ServeJSON()
}

// @Title LoadConfig
// @Description Locad ARPConfig
// @router /config [put]
func (n *HwaddmngController) LoadConfig() {
	// errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    // if errToken != nil {
    //     n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
    //     n.ServeJSON()
    //     return
    // }    
    // permissions := []string{"AddMacIp"}
    // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    // if permissionsErr != nil || hasPermission == false {
    //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    // }else{         
        anode := make(map[string]string)
        // json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "hwaddmng"
        anode["router"] = "@router / [put]"
        
        models.LoadConfig(anode)
        
        n.Data["json"] = pcap.ArpmainReturn
    // }
    n.ServeJSON()
}

// @Title Config
// @Description Config ARPConfig
// @router /config [post]
func (n *HwaddmngController) Config() {
	// errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    // if errToken != nil {
    //     n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
    //     n.ServeJSON()
    //     return
    // }    
    // permissions := []string{"AddMacIp"}
    // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    // if permissionsErr != nil || hasPermission == false {
    //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    // }else{         
        anode := make(map[string]string)
        anodeIface := make(map[string]interface{})

        json.Unmarshal(n.Ctx.Input.RequestBody, &anodeIface)

        anode["action"] = "POST"
        anode["controller"] = "hwaddmng"
        anode["router"] = "@router / [post]"
        
        models.Config(anodeIface,anode)
        
        
        n.Data["json"] = pcap.ArpmainReturn
    // }
    n.ServeJSON()
}

// @Title DbManagement
// @Description DB Management
// @router /db [post]
func (n *HwaddmngController) Db() {
	// errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    // if errToken != nil {
    //     n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
    //     n.ServeJSON()
    //     return
    // }    
    // permissions := []string{"AddMacIp"}
    // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    // if permissionsErr != nil || hasPermission == false {
    //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    // }else{         
        anode := make(map[string]string)
        anodeIface := make(map[string]string)

        json.Unmarshal(n.Ctx.Input.RequestBody, &anodeIface)

        anode["action"] = "POST"
        anode["controller"] = "hwaddmng"
        anode["router"] = "@router / [post]"
        
        models.Db(anodeIface,anode)

    // }
    // n.ServeJSON()
}

// @Title Config
// @Description Config ARPConfig
// @router /config [get]
func (n *HwaddmngController) GetConfig() {
	// errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    // if errToken != nil {
    //     n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
    //     n.ServeJSON()
    //     return
    // }    
    // permissions := []string{"AddMacIp"}
    // hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)    
    // if permissionsErr != nil || hasPermission == false {
    //     n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    // }else{         
        anode := make(map[string]string)
        anodeIface := make(map[string]interface{})

        anode["action"] = "GET"
        anode["controller"] = "hwaddmng"
        anode["router"] = "@router / [get]"
        
        models.Config(anodeIface,anode)
        
        
        n.Data["json"] = pcap.ArpmainReturn
    // }
    n.ServeJSON()
}