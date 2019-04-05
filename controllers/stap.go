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
	logs.Info ("stap controller -> GetAllServers")

	servers, err := models.GetAllServers()

	n.Data["json"] = servers

	if err != nil {
        logs.Info("GetAllServers JSON RECEIVED -- ERROR : %s", err.Error())
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
// @router /ping/:uuid [get]
func (n *StapController) PingStap() {
    logs.Info ("Stap controller -> PingStap")
    uuid := n.GetString(":uuid")
    logs.Info("Ping Stap uuid = "+uuid)
	server,err := models.PingStap(uuid)
	n.Data["json"] = server
	if err != nil {
        logs.Info("PingStap ERROR: %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    n.ServeJSON()
}

// @Title RunStap
// @Description Run Stap system
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /RunStap/:uuid [put]
func (n *StapController) RunStap() {
    logs.Info("RunStap -> In")
    uuid := n.GetString(":uuid")
    data,err := models.RunStap(uuid)
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
// @router /StopStap/:uuid [put]
func (n *StapController) StopStap() {
    logs.Info("StopStap -> In")
    uuid := n.GetString(":uuid")
    data,err := models.StopStap(uuid)
    n.Data["json"] = data
    if err != nil {
        logs.Info("StopStap OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("StopStap -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title RunStapServer
// @Description Run specific Stap server
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /RunStapServer/:serveruuid [put]
func (n *StapController) RunStapServer() {
    logs.Info("RunStapServer -> In")
    serveruuid := n.GetString(":serveruuid")
    data,err := models.RunStapServer(serveruuid)
    n.Data["json"] = data
    //logs.Warn("data RunStapServer -->"+data)
    if err != nil {
        logs.Info("RunStapServer OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("RunStapServer -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title StopStapServer
// @Description Stop specific Stap server
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /StopStapServer/:serveruuid [put]
func (n *StapController) StopStapServer() {
    logs.Info("StopStapServer -> In")
    serveruuid := n.GetString(":serveruuid")
    data,err := models.StopStapServer(serveruuid)
    n.Data["json"] = data
    if err != nil {
        logs.Info("StopStapServer OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("StopStapServer -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title PingServerStap
// @Description ping stap servers
// @Success 200 {object} models.stap
// @router /PingServerStap/:server [get]
func (n *StapController) PingServerStap() {
    logs.Info ("stap controller -> PingServerStap")
    server := n.GetString(":server")
    logs.Info("Ping Stap server = "+server)
	data, err := models.PingServerStap(server)
	n.Data["json"] = data
	if err != nil {
        logs.Info("PingServerStap OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title DeleteStapServer
// @Description Run specific Stap server
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /DeleteStapServer/:serveruuid [put]
func (n *StapController) DeleteStapServer() {
    logs.Info("DeleteStapServer -> In")
    serveruuid := n.GetString(":serveruuid")
    data,err := models.DeleteStapServer(serveruuid)
    n.Data["json"] = data
    if err != nil {
        logs.Info("DeleteStapServer OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("DeleteStapServer -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}