package controllers

import (
    "owlhnode/models"
    "github.com/astaxie/beego"
    "owlhnode/validation"
)

type ChangecontrolController struct {
    beego.Controller
}

// @Title GetChangeControlNode
// @Description Get changeControl database values
// @Success 200 {object} models.changecontrol
// @router / [get]
func (n *ChangecontrolController) GetChangeControlNode() {  
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token":"none"}
        n.ServeJSON()
        return
    }    
    permissions := []string{"GetChangeControlNode"}
    hasPermission,permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("uuid"), "any", permissions)    
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false","permissions":"none"}
    }else{    
        data, err := models.GetChangeControlNode()
        n.Data["json"] = data
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }  
    n.ServeJSON()
}