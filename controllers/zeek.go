package controllers

import (
    "owlhnode/models"
    "owlhnode/zeek"

    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "encoding/json"
)

type ZeekController struct {
    beego.Controller
}

// @Title GetZeek
// @Description get Zeek status
// @Success 200 {object} models.zeek
// @router / [get]
func (m *ZeekController) Get() {
    logs.Info ("Zeek controller -> GET")
    mstatus,err := models.GetZeek()
    m.Data["json"] = mstatus
    if err != nil {
        logs.Info("GetZeek OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}

// @Title Set
// @Description Update Zeek configuration as well as manage it
// @Description This should answer with a GetZeek so origin can verify status 
// @Success 200 {object} models.zeek
// @router / [put]
func (m *ZeekController) Set() {
    logs.Info ("Zeek controller -> PUT")
    var zeekdata zeek.Zeek
    json.Unmarshal(m.Ctx.Input.RequestBody, &zeekdata)

    zeekdata.Extra= map[string]string{}
    zeekdata.Extra["action"] = "PUT"
    zeekdata.Extra["controller"] = "ZEEK"
    zeekdata.Extra["router"] = "@router / [put]"

    mstatus,err := models.SetZeek(zeekdata)
    m.Data["json"] = mstatus
    if err != nil {
        logs.Info("Set Zeek OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}


// @Title RunZeek
// @Description Run zeek system
// @Success 200 {object} models.zeek
// @Failure 403 body is empty
// @router /RunZeek [put]
func (n *ZeekController) RunZeek() {
    logs.Info("RunZeek -> In")
    anode := map[string]string{}
    anode["action"] = "PUT"
    anode["controller"] = "ZEEK"
    anode["router"] = "@router /RunZeek [put]"
    logs.Info("============")
    logs.Info("ZEEK - RunZeek")
    for key :=range anode {
        logs.Info(key +" -> "+anode[key])
    }
    data,err := models.RunZeek()
    n.Data["json"] = data
    if err != nil {
        logs.Info("RunZeek OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("RunZeek -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title StopZeek
// @Description Run zeek system
// @Success 200 {object} models.zeek
// @Failure 403 body is empty
// @router /StopZeek [put]
func (n *ZeekController) StopZeek() {
    logs.Info("StopZeek -> In")
    var anode map[string]string
    anode["action"] = "PUT"
    anode["controller"] = "ZEEK"
    anode["router"] = "@router /StopZeek [put]"
    logs.Info("============")
    logs.Info("ZEEK - StopZeek")
    for key :=range anode {
        logs.Info(key +" -> "+anode[key])
    }
    data,err := models.StopZeek()
    n.Data["json"] = data
    if err != nil {
        logs.Info("StopZeek OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("StopZeek -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title ChangeZeekMode
// @Description get Zeek status
// @Success 200 {object} models.zeek
// @router /changeZeekMode [put]
func (m *ZeekController) ChangeZeekMode() {
    var anode map[string]string
    json.Unmarshal(m.Ctx.Input.RequestBody, &anode)

    anode["action"] = "PUT"
    anode["controller"] = "ZEEK"
    anode["router"] = "@router /changeZeekMode [put]"
    err := models.ChangeZeekMode(anode)
    m.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}

// @Title AddClusterValue
// @Description get Zeek status
// @Success 200 {object} models.zeek
// @router /addClusterValue [post]
func (m *ZeekController) AddClusterValue() {
    var anode map[string]string
    json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "ZEEK"
    anode["router"] = "@router /addClusterValue [post]"

    err := models.AddClusterValue(anode)
    m.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}

// @Title PingCluster
// @Description get Zeek cluster elements
// @Success 200 {object} models.zeek
// @router /pingCluster [get]
func (m *ZeekController) PingCluster() {
    data,err := models.PingCluster()
    m.Data["json"] = data
    if err != nil {
        m.Data["json"] = map[string]map[string]string{"error":{"ack": "false", "error": err.Error()}}
    }
    m.ServeJSON()
}

// @Title EditClusterValue
// @Description Edit Zeek status
// @Success 200 {object} models.zeek
// @router /editClusterValue [put]
func (m *ZeekController) EditClusterValue() {
    var anode map[string]string
    json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "ZEEK"
    anode["router"] = "@router /editClusterValue [put]"
    err := models.EditClusterValue(anode)
    m.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}

// @Title DeleteClusterValue
// @Description Delete Zeek status
// @Success 200 {object} models.zeek
// @router /deleteClusterValue [delete]
func (m *ZeekController) DeleteClusterValue() {
    var anode map[string]string
    json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
    anode["action"] = "DELETE"
    anode["controller"] = "ZEEK"
    anode["router"] = "@router /deleteClusterValue [delete]"
    err := models.DeleteClusterValue(anode)
    m.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}

// @Title SyncCluster
// @Description Sync Zeek cluster
// @Success 200 {object} models.zeek
// @router /syncCluster [put]
func (m *ZeekController) SyncCluster() {
    anode := map[string]string{}
    // var anode map[string]string
    json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
    anode["action"] = "PUT"
    anode["controller"] = "ZEEK"
    anode["router"] = "@router /syncCluster [put]"
    err := models.SyncCluster(anode)
    m.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    m.ServeJSON()
}

// @Title SaveConfigFile
// @Description Save Configuration files from Master
// @Success 200 {object} models.Node
// @Failure 403 body is empty
// @router / [post]
func (n *SuricataController) SaveConfigFile() {
    var anode map[string][]byte
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

    logs.Info("ACTION -> POST")
    logs.Info("CONTROLLER -> ZEEK")
    logs.Info("ROUTER -> @router / [post]")
    for key := range anode {
        logs.Info("key -> "+key)
    }

    err := models.SaveConfigFile(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        logs.Info("Save configuration files -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}
