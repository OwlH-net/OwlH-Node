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
// @Description Add a new server to stap
// @Success 200 {object} models.stap
// @router / [post]
func (n *StapController) AddServer() {
	logs.Info ("stap controller -> AddServer")
	
	var newServer map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &newServer)

	err := models.AddServer(newServer)

	n.Data["json"] = map[string]string{"ack": "true"}

	if err != nil {
        logs.Info("AddServer JSON RECEIVED -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}

    n.ServeJSON()
}

// @Title GetAllServers
// @Description get all servers stap
// @Success 200 {object} models.stap
// @router / [get]
func (n *StapController) GetAllServers() {
	logs.Info ("stap controller -> AddServer")

	servers, err := models.GetAllServers()

	n.Data["json"] = servers

	if err != nil {
        logs.Info("AddServer JSON RECEIVED -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    n.ServeJSON()
}