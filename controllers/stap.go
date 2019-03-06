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

// @Title GetServer
// @Description get a server stap
// @Success 200 {object} models.stap
// @router /server/:uuid [get]
func (n *StapController) GetServer() {
	logs.Info ("stap controller -> GetServer")
	uuid := n.GetString(":uuid") 
	server, err := models.GetServer(uuid)

	n.Data["json"] = server

	if err != nil {
        logs.Info("GetServer JSON RECEIVED -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    n.ServeJSON()
}

// @Title PingStap
// @Description ping stap servers
// @Success 200 {object} models.stap
// @router /ping [get]
func (n *StapController) PingStap() {
	logs.Info ("stap controller -> PingStap")
	server, err := models.PingStap()

	n.Data["json"] = server

	if err != nil {
        logs.Info("PingStap JSON RECEIVED -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    n.ServeJSON()
}

// @Title RunStap
// @Description Run Stap system
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /RunStap [put]
func (n *StapController) RunStap() {
    logs.Info("RunStap -> In")
    data,err := models.RunStap()
    n.Data["json"] = data
    if err != nil {
        logs.Info("RunStap OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("RunStap -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title StopStap
// @Description Run Stap system
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /StopStap [put]
func (n *StapController) StopStap() {
    logs.Info("StopStap -> In")
    data,err := models.StopStap()
    n.Data["json"] = data
    if err != nil {
        logs.Info("StopStap OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("StopStap -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}