package controllers

import (
	"owlhnode/models"
	"encoding/json"
	// "strconv"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
)

type StapController struct {
	beego.Controller
}

// @Title AddServer
// @Description get Surucata status
// @Success 200 {object} models.stap
// @router / [post]
func (n *StapController) AddServer() {
	logs.Info ("stap controller -> AddServer")
	
	var newServer map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &newServer)

	err := models.AddServer(newServer)

	n.Data["json"] = map[string]string{"ack": "true"}

	if err != nil {
        logs.Info("BPF JSON RECEIVED -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}

    n.ServeJSON()
}