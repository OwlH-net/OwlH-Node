package controllers

import (
    "encoding/json"
    "github.com/astaxie/beego"
    // "github.com/astaxie/beego/logs"
    "owlhnode/models"
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
        err := models.AddMacIp(anode)
        anode["action"] = "POST"
        anode["controller"] = "hwaddmng"
        anode["router"] = "@router / [post]"
    
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    // }
    n.ServeJSON()
}