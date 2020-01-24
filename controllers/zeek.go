package controllers

import (
    "owlhnode/models"
    "owlhnode/zeek"
    "owlhnode/validation"
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
func (n *ZeekController) Get() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        logs.Info ("Zeek controller -> GET")
        mstatus,err := models.GetZeek()
        n.Data["json"] = mstatus
        if err != nil {
            logs.Info("GetZeek OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title Set
// @Description Update Zeek configuration as well as manage it
// @Description This should answer with a GetZeek so origin can verify status 
// @Success 200 {object} models.zeek
// @router / [put]
func (n *ZeekController) Set() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        logs.Info ("Zeek controller -> PUT")
        var zeekdata zeek.Zeek
        json.Unmarshal(n.Ctx.Input.RequestBody, &zeekdata)
    
        zeekdata.Extra= map[string]string{}
        zeekdata.Extra["action"] = "PUT"
        zeekdata.Extra["controller"] = "ZEEK"
        zeekdata.Extra["router"] = "@router / [put]"
    
        mstatus,err := models.SetZeek(zeekdata)
        n.Data["json"] = mstatus
        if err != nil {
            logs.Info("Set Zeek OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}


// @Title RunZeek
// @Description Run zeek system
// @Success 200 {object} models.zeek
// @Failure 403 body is empty
// @router /RunZeek [put]
func (n *ZeekController) RunZeek() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
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
    }
    n.ServeJSON()
}

// @Title StopZeek
// @Description Run zeek system
// @Success 200 {object} models.zeek
// @Failure 403 body is empty
// @router /StopZeek [put]
func (n *ZeekController) StopZeek() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
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
    }
    n.ServeJSON()
}

// @Title ChangeZeekMode
// @Description get Zeek status
// @Success 200 {object} models.zeek
// @router /changeZeekMode [put]
func (n *ZeekController) ChangeZeekMode() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    
        anode["action"] = "PUT"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /changeZeekMode [put]"
        err := models.ChangeZeekMode(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title AddClusterValue
// @Description get Zeek status
// @Success 200 {object} models.zeek
// @router /addClusterValue [post]
func (n *ZeekController) AddClusterValue() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /addClusterValue [post]"
    
        err := models.AddClusterValue(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title PingCluster
// @Description get Zeek cluster elements
// @Success 200 {object} models.zeek
// @router /pingCluster [get]
func (n *ZeekController) PingCluster() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        data,err := models.PingCluster()
        n.Data["json"] = data
        if err != nil {
            n.Data["json"] = map[string]map[string]string{"error":{"ack": "false", "error": err.Error()}}
        }
    }
    n.ServeJSON()
}

// @Title EditClusterValue
// @Description Edit Zeek status
// @Success 200 {object} models.zeek
// @router /editClusterValue [put]
func (n *ZeekController) EditClusterValue() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /editClusterValue [put]"
        err := models.EditClusterValue(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title DeleteClusterValue
// @Description Delete Zeek status
// @Success 200 {object} models.zeek
// @router /deleteClusterValue [delete]
func (n *ZeekController) DeleteClusterValue() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "DELETE"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /deleteClusterValue [delete]"
        err := models.DeleteClusterValue(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title SyncCluster
// @Description Sync Zeek cluster
// @Success 200 {object} models.zeek
// @router /syncCluster [put]
func (n *ZeekController) SyncCluster() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        anode := map[string]string{}
        // var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /syncCluster [put]"
        err := models.SyncCluster(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title SavePolicyFiles
// @Description Save Configuration files from Master
// @Success 200 {object} models.Node
// @Failure 403 body is empty
// @router / [post]
func (n *ZeekController) SavePolicyFiles() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        var anode map[string]map[string][]byte
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    
        logs.Info("ACTION -> POST")
        logs.Info("CONTROLLER -> ZEEK")
        logs.Info("ROUTER -> @router / [post]")
        for key := range anode {
            logs.Info("key -> "+key)
        }
    
        err := models.SavePolicyFiles(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Info("Save configuration files -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title SyncClusterFile
// @Description Sync Zeek cluster file
// @Success 200 {object} models.zeek
// @router /syncClusterFile [put]
func (n *ZeekController) SyncClusterFile() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        anode := map[string][]byte{}
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    
        logs.Info("ACTION -> POST")
        logs.Info("CONTROLLER -> ZEEK")
        logs.Info("ROUTER -> @router / [post]")
        for key := range anode {
            logs.Info("key -> "+key)
        }
        
        err := models.SyncClusterFile(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title LaunchZeekMainConf
// @Description Sync Zeek cluster file
// @Success 200 {object} models.zeek
// @router /LaunchZeekMainConf [put]
func (n *ZeekController) LaunchZeekMainConf() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        anode := map[string]string{}
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    
        logs.Info("ACTION -> POST")
        logs.Info("CONTROLLER -> ZEEK")
        logs.Info("ROUTER -> @router /LaunchZeekMainConf [put]")
        for key := range anode {
            logs.Info("key -> "+key)
        }
        
        err := models.LaunchZeekMainConf(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// // @Title SaveZeekValues
// // @Description Edit Zeek expert values
// // @Success 200 {object} models.zeek
// // @router /saveZeekValues [put]
// func (n *ZeekController) SaveZeekValues() {
//     var anode map[string]string
//     json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
//     anode["action"] = "PUT"
//     anode["controller"] = "ZEEK"
//     anode["router"] = "@router /SaveZeekValues [put]"
//     err := models.SaveZeekValues(anode)
//     n.Data["json"] = map[string]string{"ack": "true"}
//     if err != nil {
//         n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
//     }
//     n.ServeJSON()
// }

// @Title SyncZeekValues
// @Description Sync Zeek cluster file
// @Success 200 {object} models.zeek
// @router /syncZeekValues [put]
func (n *ZeekController) SyncZeekValues() {
    err := validation.CheckToken(n.Ctx.Input.Header("token"))
    if err != nil {
        logs.Error("Error validating token from master")
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else{         
        anode := map[string]string{}
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    
        logs.Info("ACTION -> POST")
        logs.Info("CONTROLLER -> ZEEK")
        logs.Info("ROUTER -> @router /SyncZeekValues [put]")
        for key := range anode {
            logs.Info("key -> "+key)
        }
        
        err := models.SyncZeekValues(anode)
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}