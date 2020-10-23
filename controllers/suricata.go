package controllers

import (
    "encoding/json"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "owlhnode/models"
    "owlhnode/validation"
)

type SuricataController struct {
    beego.Controller
}

// @Title GetSuricata
// @Description get Surucata status
// @Success 200 {object} models.suricata
// @router / [get]
func (n *SuricataController) Get() {
    mstatus, err := models.GetSuricata()

    n.Data["json"] = mstatus
    if err != nil {
        logs.Info("Suricata controller -> GET -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// // @Title Get Suricata BPF
// // @Description get Suricata BPF from filter.bpf file
// // @Success 200 {object} models.suricata
// // @router /bpf [get]
// func (n *SuricataController) GetBPF() {
//     bpf,err := models.GetBPF()
//     n.Data["json"] = bpf
//     if err != nil {
//         logs.Info("GetBPF OUT -- ERROR : %s", err.Error())
//         n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
//     }
//     n.ServeJSON()
// }

// @Title PUT Suricata BPF
// @Description Set Suricata BPF into filter.bpf file
// @Success 200 {object} models.suricata
// @router /bpf [put]
func (n *SuricataController) SetBPF() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"SetBPF"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        anode["action"] = "PUT"
        anode["controller"] = "SURICATA"
        anode["router"] = "@router /bpf [put]"

        err := models.SetBPF(anode, n.Ctx.Input.Header("user"))

        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Info("BPF JSON RECEIVED -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title SyncRulesetFromMaster
// @Description Sync ruleset file from master
// @Success 200 {object} models.Node
// @Failure 403 body is empty
// @router /sync [put]
func (n *SuricataController) SyncRulesetFromMaster() {
    // permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    // if err != nil {
    //     logs.Error("Suricata Error validating token from master")
    //     logs.Error(err.Error())
    //     n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    // }else if !permissions{
    //     n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    // }else{
    var anode map[string][]byte
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

    logs.Info("ACTION -> PUT")
    logs.Info("CONTROLLER -> SURICATA")
    logs.Info("ROUTER -> @router /sync [put]")
    for key := range anode {
        logs.Info("key -> " + key)
    }

    err := models.SyncRulesetFromMaster(anode, n.Ctx.Input.Header("user"))

    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        logs.Info("Ruleset retrieve OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    // }
    n.ServeJSON()
}

// @Title SaveConfigFile
// @Description Save Configuration files from Master
// @Success 200 {object} models.Node
// @Failure 403 body is empty
// @router / [post]
func (n *SuricataController) SaveConfigFile() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"SaveConfigFile"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string][]byte
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        logs.Info("ACTION -> POST")
        logs.Info("CONTROLLER -> SURICATA")
        logs.Info("ROUTER -> @router / [post]")
        for key := range anode {
            logs.Info("key -> " + key)
        }

        err := models.SaveConfigFile(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Info("Save configuration files -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title RunSuricata
// @Description Run suricata system
// @Success 200 {object} models.suricata
// @Failure 403 body is empty
// @router /RunSuricata [put]
func (n *SuricataController) RunSuricata() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"RunSuricata"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        anode["action"] = "PUT"
        anode["controller"] = "SURICATA"
        anode["router"] = "@router /RunSuricata [put]"
        logs.Info("============")
        logs.Info("SURICATA - RunSuricata")
        for key := range anode {
            logs.Info(key + " -> " + anode[key])
        }
        data, err := models.RunSuricata(n.Ctx.Input.Header("user"))
        n.Data["json"] = data
        if err != nil {
            logs.Info("RunSuricata OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title SuricataConfigurationTest
// @Description Check Configuration
// @Success 200 {object} models.suricata
// @router /check [put]
func (n *SuricataController) SuricataConfigurationTest() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"RunSuricata"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode = map[string]string{}
        anode["action"] = "PUT"
        anode["controller"] = "SURICATA"
        anode["router"] = "@router /RunSuricata [put]"
        logs.Info("============")
        logs.Info("SURICATA - SuricataConfigurationTest")
        for key := range anode {
            logs.Info(key + " -> " + anode[key])
        }
        data, err := models.SuricataConfigurationTest(n.Ctx.Input.Header("uuid"))
        n.Data["json"] = data
        if err != nil {
            logs.Info("SuricataConfigurationTest OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title DumpSuricataConfiguration
// @Description Dump Suricata Configuration
// @Success 200 {object} models.suricata
// @router /dump [get]
func (n *SuricataController) SuricataDumpCurrentConfig() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"RunSuricata"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode = map[string]string{}
        anode["action"] = "GET"
        anode["controller"] = "SURICATA"
        anode["router"] = "@router /dump [get]"
        logs.Info("============")
        logs.Info("SURICATA - DumpSuricataConfiguration")
        for key := range anode {
            logs.Info(key + " -> " + anode[key])
        }
        data, err := models.SuricataDumpCurrentConfig()
        n.Data["json"] = data
        if err != nil {
            logs.Info("DumpSuricataConfiguration OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title StopSuricata
// @Description Run suricata system
// @Success 200 {object} models.suricata
// @Failure 403 body is empty
// @router /StopSuricata [put]
func (n *SuricataController) StopSuricata() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"StopSuricata"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode = map[string]string{}
        anode["action"] = "PUT"
        anode["controller"] = "SURICATA"
        anode["router"] = "@router /StopSuricata [put]"
        logs.Info("============")
        logs.Info("SURICATA - StopSuricata")
        for key := range anode {
            logs.Info(key + " -> " + anode[key])
        }
        data, err := models.StopSuricata(n.Ctx.Input.Header("user"))
        n.Data["json"] = data
        if err != nil {
            logs.Info("StopSuricata OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title GetSuricataServices
// @Description get all Suricata services
// @Success 200 {object} models.suricata
// @router /get [get]
func (n *SuricataController) GetSuricataServices() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"GetSuricataServices"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        servicesSuricata, err := models.GetSuricataServices(n.Ctx.Input.Header("user"))
        n.Data["json"] = servicesSuricata

        if err != nil {
            logs.Info("GetSuricataServices ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title StartSuricataMainConf
// @Description Start Suricata Main Conf
// @Success 200 {object} models.suricata
// @Failure 403 body is empty
// @router /StartSuricataMainConf [put]
func (n *SuricataController) StartSuricataMainConf() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"StartSuricataMainConf"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        logs.Info("ACTION -> PUT")
        logs.Info("CONTROLLER -> SURICATA")
        logs.Info("ROUTER -> @router /StartSuricataMainConf [put]")
        for key := range anode {
            logs.Info("key -> " + key)
        }

        err := models.StartSuricataMainConf(anode, n.Ctx.Input.Header("user"))
        if err != nil {
            logs.Info("StartSuricataMainConf OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title StopSuricataMainConf
// @Description Stop Suricata Main Conf
// @Success 200 {object} models.suricata
// @Failure 403 body is empty
// @router /StopSuricataMainConf [put]
func (n *SuricataController) StopSuricataMainConf() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"StopSuricataMainConf"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        logs.Info("ACTION -> PUT")
        logs.Info("CONTROLLER -> SURICATA")
        logs.Info("ROUTER -> @router /StopSuricataMainConf [put]")
        for key := range anode {
            logs.Info("key -> " + key)
        }

        err := models.StopSuricataMainConf(anode, n.Ctx.Input.Header("user"))
        if err != nil {
            logs.Info("StopSuricataMainConf OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title KillSuricataMainConf
// @Description Kill Suricata Main Conf
// @Success 200 {object} models.suricata
// @Failure 403 body is empty
// @router /KillSuricataMainConf [put]
func (n *SuricataController) KillSuricataMainConf() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"KillSuricataMainConf"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        logs.Info("ACTION -> PUT")
        logs.Info("CONTROLLER -> SURICATA")
        logs.Info("ROUTER -> @router /KillSuricataMainConf [put]")
        for key := range anode {
            logs.Info("key -> " + key)
        }

        err := models.KillSuricataMainConf(anode, n.Ctx.Input.Header("user"))
        if err != nil {
            logs.Info("KillSuricataMainConf OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title ReloadSuricataMainConf
// @Description Reload Suricata Main Conf
// @Success 200 {object} models.suricata
// @Failure 403 body is empty
// @router /ReloadSuricataMainConf [put]
func (n *SuricataController) ReloadSuricataMainConf() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"ReloadSuricataMainConf"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        logs.Info("ACTION -> PUT")
        logs.Info("CONTROLLER -> SURICATA")
        logs.Info("ROUTER -> @router /ReloadSuricataMainConf [put]")
        for key := range anode {
            logs.Info("key -> " + key)
        }

        err := models.ReloadSuricataMainConf(anode, n.Ctx.Input.Header("user"))
        if err != nil {
            logs.Info("ReloadSuricataMainConf OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title GetMD5files
// @Description Reload Suricata Main Conf
// @Success 200 {object} models.suricata
// @Failure 403 body is empty
// @router /getMD5files [put]
func (n *SuricataController) GetMD5files() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        var errorResponse = map[string]map[string]string{}
        errorResponse["hasError"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.Data["json"] = errorResponse
        n.ServeJSON()
        return
    }
    permissions := []string{"GetMD5files"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        var errorResponse = map[string]map[string]string{}
        errorResponse["hasError"] = map[string]string{"ack": "false", "permissions": "none", "error": "Not enough permissions"}
        n.Data["json"] = errorResponse
    } else {
        var anode map[string]map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        logs.Info("ACTION -> PUT")
        logs.Info("CONTROLLER -> SURICATA")
        logs.Info("ROUTER -> @router /GetMD5files [put]")
        for key := range anode {
            logs.Info("key -> " + key)
        }

        data, err := models.GetMD5files(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = data
        if err != nil {
            logs.Info("GetMD5files OUT -- ERROR : %s", err.Error())
            var errorResponse = map[string]map[string]string{}
            errorResponse["hasError"] = map[string]string{"ack": "false", "error": err.Error()}
            n.Data["json"] = errorResponse
        }
    }
    n.ServeJSON()
}

// @Title AddSuricataService()
// @Description Add new Suricata Service
// @Success 200 {object} models.suricata
// @router /service [put]
func (n *SuricataController) AddSuricataService() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"AddPluginService"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "SURICATA"
        anode["router"] = "@router /service [put]"
        err := models.AddSuricataService(anode, n.Ctx.Input.Header("user"))

        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title GetSuricataRulesets()
// @Description Add new Suricata Service
// @Success 200 {object} models.suricata
// @router /getSuricataRulesets [get]
func (n *SuricataController) GetSuricataRulesets() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"GetSuricataRulesets"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        data,err := models.GetSuricataRulesets()

        n.Data["json"] = data
        if err != nil {
            var errorResponse = map[string]map[string]string{}
            errorResponse["hasError"] = map[string]string{"ack": "false", "error": err.Error()}
            n.Data["json"] = errorResponse
        }
    }
    n.ServeJSON()
}