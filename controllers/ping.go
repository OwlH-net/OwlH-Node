package controllers

import (
    "github.com/astaxie/beego"
    "owlhnode/models"
    "encoding/json"
//    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)

type PingController struct {
    beego.Controller
}

// @Title PingNode
// @Description get ping for node
// @router / [get]
func (n *PingController) PingNode() {
    n.Data["json"] = map[string]string{"ack": "true"}
    n.ServeJSON()
}

// @Title UpdateNodeData
// @Description update node data
// @router /updateNode [put]
func (n *PingController) UpdateNodeData() {
    var anode map[string]map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    logs.Info("ACTION -> PUT")
    logs.Info("CONTROLLER -> PING")
    logs.Info("ROUTER -> @router /updateNode [put]")
    for key := range anode {
        logs.Info("key -> "+key)
    }
    err := models.UpdateNodeData(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title PingService
// @Description get ping for node
// @router /services [get]
func (n *PingController) PingService() {
    err := models.PingService()
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title DeployService
// @Description get ping for node
// @router /deployservice [put]
func (n *PingController) DeployService() {
    var anode map[string]map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    logs.Info("ACTION -> PUT")
    logs.Info("CONTROLLER -> PING")
    logs.Info("ROUTER -> @router /deployservice [put]")
    err := models.DeployService(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title GetMainconfData
// @Description get ping for node
// @router /mainconf [get]
func (n *PingController) GetMainconfData() {
    data,err := models.GetMainconfData()
    n.Data["json"] = data
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title PingPluginsNode
// @Description PingPluginsNode status
// @Success 200 {object} models.ports
// @router /PingPluginsNode [get]
func (m *PingController) PingPluginsNode() {
    data, err := models.PingPluginsNode()
    m.Data["json"] = data
    if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}