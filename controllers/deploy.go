package controllers

import (
	"owlhnode/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type DeployController struct {
	beego.Controller
}

// @Title Deploy
// @Description save changes over requested file on webpage "edit.html"
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router / [put]
func (n *DeployController) Deploy() {
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    err := models.Deploy(anode)

    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title Deploy
// @Description save changes over requested file on webpage "edit.html"
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router / [get]
func (n *DeployController) CheckDeployFiles() {
    anode := models.CheckDeployFiles()
    n.Data["json"] = anode
    n.ServeJSON()
}