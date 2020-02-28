package controllers

import (
    "encoding/json"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "owlhnode/models"
    "owlhnode/validation"
)

type StapController struct {
    beego.Controller
}

// @Title AddServer
// @Description Add a new server to stap
// @Success 200 {object} models.stap
// @router / [post]
func (n *StapController) AddServer() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "post")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        logs.Info ("stap controller -> AddServer")
        
        var newServer map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &newServer)
        newServer["action"] = "POST"
        newServer["controller"] = "STAP"
        newServer["router"] = "@router / [post]"
        err := models.AddServer(newServer)
    
        n.Data["json"] = map[string]string{"ack": "true"}
    
        if err != nil {
            logs.Info("AddServer JSON RECEIVED -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title GetAllServers
// @Description get all servers stap
// @Success 200 {object} models.stap
// @router / [get]
func (n *StapController) GetAllServers() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        logs.Info ("stap controller -> GetAllServers")
    
        servers, err := models.GetAllServers()
    
        n.Data["json"] = servers
    
        if err != nil {
            logs.Info("GetAllServers JSON RECEIVED -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title GetServer
// @Description get a server stap
// @Success 200 {object} models.stap
// @router /server/:uuid [get]
func (n *StapController) GetServer() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        logs.Info ("stap controller -> GetServer")
        uuid := n.GetString(":uuid") 
        server, err := models.GetServer(uuid)
    
        n.Data["json"] = server
    
        if err != nil {
            logs.Info("GetServer JSON RECEIVED -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title PingStap
// @Description ping stap servers
// @Success 200 {object} models.stap
// @router /ping/:uuid [get]
func (n *StapController) PingStap() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        logs.Info ("Stap controller -> PingStap")
        uuid := n.GetString(":uuid")
        logs.Info("Ping Stap uuid = "+uuid)
        server,err := models.PingStap(uuid)
        n.Data["json"] = server
        if err != nil {
            logs.Info("PingStap ERROR: %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title RunStap
// @Description Run Stap system
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /RunStap/:uuid [put]
func (n *StapController) RunStap() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        logs.Info("RunStap -> In")
        uuid := n.GetString(":uuid")
        var anode map[string]string
        anode["action"] = "PUT"
        anode["controller"] = "STAP"
        anode["router"] = "@router /RunStap/:uuid [put]"
        anode["uuid"] = uuid
        logs.Info("============")
        logs.Info("STAP - RunStap")
        for key :=range anode {
            logs.Info(key +" -> "+anode[key])
        }
    
        data,err := models.RunStap(uuid)
        n.Data["json"] = data
        if err != nil {
            logs.Info("RunStap OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        logs.Info("RunStap -> OUT -> %s", n.Data["json"])
    }
    n.ServeJSON()
}

// @Title StopStap
// @Description Run Stap system
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /StopStap/:uuid [put]
func (n *StapController) StopStap() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        logs.Info("StopStap -> In")
        uuid := n.GetString(":uuid")
        var anode map[string]string
        anode["action"] = "PUT"
        anode["controller"] = "STAP"
        anode["router"] = "@router /StopStap/:uuid [put]"
        anode["uuid"] = uuid
        logs.Info("============")
        logs.Info("STAP - RunStap")
        for key :=range anode {
            logs.Info(key +" -> "+anode[key])
        }
        data,err := models.StopStap(uuid)
        n.Data["json"] = data
        if err != nil {
            logs.Info("StopStap OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        logs.Info("StopStap -> OUT -> %s", n.Data["json"])
    }
    n.ServeJSON()
}

// @Title RunStapServer
// @Description Run specific Stap server
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /RunStapServer/:serveruuid [put]
func (n *StapController) RunStapServer() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        logs.Info("RunStapServer -> In")
        serveruuid := n.GetString(":serveruuid")
            var anode map[string]string
        anode["action"] = "PUT"
        anode["controller"] = "STAP"
        anode["router"] = "@router /RunStapServer/:serveruuid [put]"
        anode["serveruuid"] = serveruuid
        logs.Info("============")
        logs.Info("STAP - RunStap")
        for key :=range anode {
            logs.Info(key +" -> "+anode[key])
        }
        data,err := models.RunStapServer(serveruuid)
        n.Data["json"] = data
        //logs.Warn("data RunStapServer -->"+data)
        if err != nil {
            logs.Info("RunStapServer OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        logs.Info("RunStapServer -> OUT -> %s", n.Data["json"])
    }
    n.ServeJSON()
}

// @Title StopStapServer
// @Description Stop specific Stap server
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /StopStapServer/:serveruuid [put]
func (n *StapController) StopStapServer() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        logs.Info("StopStapServer -> In")
        serveruuid := n.GetString(":serveruuid")
        var anode map[string]string
        anode["action"] = "PUT"
        anode["controller"] = "STAP"
        anode["router"] = "@router /StopStapServer/:serveruuid [put]"
        anode["serveruuid"] = serveruuid
        logs.Info("============")
        logs.Info("STAP - RunStap")
        for key :=range anode {
            logs.Info(key +" -> "+anode[key])
        }
        data,err := models.StopStapServer(serveruuid)
        n.Data["json"] = data
        if err != nil {
            logs.Info("StopStapServer OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        logs.Info("StopStapServer -> OUT -> %s", n.Data["json"])
    }
    n.ServeJSON()
}

// @Title PingServerStap
// @Description ping stap servers
// @Success 200 {object} models.stap
// @router /PingServerStap/:server [get]
func (n *StapController) PingServerStap() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        logs.Info ("stap controller -> PingServerStap")
        server := n.GetString(":server")
        logs.Info("Ping Stap server = "+server)
    
        data, err := models.PingServerStap(server)
        n.Data["json"] = data
        if err != nil {
            logs.Info("PingServerStap OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title DeleteStapServer
// @Description Run specific Stap server
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /DeleteStapServer/:serveruuid [put]
func (n *StapController) DeleteStapServer() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        logs.Info("DeleteStapServer -> In")
        serveruuid := n.GetString(":serveruuid")
        var anode map[string]string
        anode["action"] = "PUT"
        anode["controller"] = "STAP"
        anode["router"] = "@router /DeleteStapServer/:serveruuid [put]"
        anode["serveruuid"] = serveruuid
        logs.Info("============")
        logs.Info("STAP - RunStap")
        for key :=range anode {
            logs.Info(key +" -> "+anode[key])
        }
        
        data,err := models.DeleteStapServer(serveruuid)
        n.Data["json"] = data
        if err != nil {
            logs.Info("DeleteStapServer OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        logs.Info("DeleteStapServer -> OUT -> %s", n.Data["json"])
    }
    n.ServeJSON()
}

// @Title EditStapServer
// @Description Edit specific Stap server
// @Success 200 {object} models.Stap
// @Failure 403 body is empty
// @router /EditStapServer [put]
func (n *StapController) EditStapServer() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{    
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "STAP"
        anode["router"] = "@router /EditStapServer [put]"
    
        err := models.EditStapServer(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Info("EditStapServer OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}