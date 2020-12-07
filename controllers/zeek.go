package controllers

import (
    "encoding/json"
    "fmt"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "owlhnode/models"
    "owlhnode/validation"
    "owlhnode/zeek"
)

type ZeekController struct {
    beego.Controller
}

// @Title GetZeek
// @Description get Zeek status
// @Success 200 {object} models.zeek
// @router / [get]
func (n *ZeekController) Get() {

    logs.Info("Zeek - Get info --- --- --- ")
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        logs.Info("Zeek - Token Error")
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"GetZeek"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)

    if permissionsErr != nil || hasPermission == false {
        logs.Error("zeek - user %s doesn't has permissions for Zeek GET INFO ", n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        logs.Info("Zeek controller -> GET")
        mstatus, err := models.GetZeek(n.Ctx.Input.Header("user"))
        n.Data["json"] = mstatus
        if err != nil {
            logs.Info("GetZeek OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    logs.Debug("Zeek - Get info - return back ")

    n.ServeJSON()
}

// @Title Set
// @Description Update Zeek configuration as well as manage it
// @Description This should answer with a GetZeek so origin can verify status
// @Success 200 {object} models.zeek
// @router / [put]
func (n *ZeekController) Set() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"Set"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        logs.Info("Zeek controller -> PUT")
        var zeekdata zeek.Zeek
        json.Unmarshal(n.Ctx.Input.RequestBody, &zeekdata)

        zeekdata.Extra = map[string]string{}
        zeekdata.Extra["action"] = "PUT"
        zeekdata.Extra["controller"] = "ZEEK"
        zeekdata.Extra["router"] = "@router / [put]"

        mstatus, err := models.SetZeek(zeekdata, n.Ctx.Input.Header("user"))
        n.Data["json"] = mstatus
        if err != nil {
            logs.Info("Set Zeek OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title Run Zeekctl command
// @Description run zeekctl command stop - start - deploy - status -diag
// @Success 200 {object} models.zeek
// @router /:command [put]
func (n *ZeekController) Command() {
    // errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    // if errToken != nil {
    //     n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
    //     n.ServeJSON()
    //     return
    // }
    // permissions := []string{"RunZeekCommand"}
    // hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    // if  permissionsErr != nil || hasPermission == false {
    if false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        cmd := n.GetString(":command")
        var iserror error
        switch cmd {
        case "start":
            data, err := models.StartZeek("start", n.Ctx.Input.Header("user"))
            n.Data["json"] = data
            if err != nil {
                iserror = err
            }
        case "stop":
            data, err := models.StopZeek(n.Ctx.Input.Header("user"))
            n.Data["json"] = data
            if err != nil {
                iserror = err
            }
        case "deploy":
            data, err := models.DeployZeek(n.Ctx.Input.Header("user"))
            n.Data["json"] = data
            if err != nil {
                iserror = err
            }
        case "diag":
            data, err := models.DiagZeek(n.Ctx.Input.Header("user"))
            n.Data["json"] = data

            if err != nil {
                iserror = err
            }
        case "status":
            data, err := models.GetZeek(n.Ctx.Input.Header("user"))
            n.Data["json"] = data
            if err != nil {
                iserror = err
            }
        default:
            str := fmt.Sprintf("ZEEK - Run command %s -- ERROR : Command Unknown", cmd)
            logs.Error(str)
            n.Data["json"] = map[string]string{"ack": "false", "error": str}
        }
        if iserror != nil {
            logs.Error("ZEEK - Run command %s -- ERROR : %s", cmd, iserror.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": iserror.Error()}
        }
    }
    logs.Debug("Zeek - run command - return back ")
    n.ServeJSON()
}

// @Title RunZeek
// @Description Run zeek system
// @Success 200 {object} models.zeek
// @Failure 403 body is empty
// @router /RunZeek [put]
func (n *ZeekController) RunZeek() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"RunZeek"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        logs.Info("RunZeek -> In")
        anode := map[string]string{}
        anode["action"] = "PUT"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /RunZeek [put]"
        logs.Info("============")
        logs.Info("ZEEK - RunZeek")
        for key := range anode {
            logs.Info(key + " -> " + anode[key])
        }
        data, err := models.RunZeek(n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"StopZeek"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        logs.Info("StopZeek -> In")
        var anode map[string]string
        anode["action"] = "PUT"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /StopZeek [put]"
        logs.Info("============")
        logs.Info("ZEEK - StopZeek")
        for key := range anode {
            logs.Info(key + " -> " + anode[key])
        }
        data, err := models.StopZeek(n.Ctx.Input.Header("user"))
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
// @Description change zeek mode
// @Success 200 {object} models.zeek
// @router /changeZeekMode [put]
func (n *ZeekController) ChangeZeekMode() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"ChangeZeekMode"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        anode["action"] = "PUT"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /changeZeekMode [put]"
        err := models.ChangeZeekMode(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}

// @Title AddClusterValue
// @Description Add zeek cluster value
// @Success 200 {object} models.zeek
// @router /addClusterValue [post]
func (n *ZeekController) AddClusterValue() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"AddClusterValue"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /addClusterValue [post]"

        err := models.AddClusterValue(anode, n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        var errorResponse = map[string]map[string]string{}
        errorResponse["hasError"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.Data["json"] = errorResponse
        n.ServeJSON()
        return
    }
    permissions := []string{"PingCluster"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        var errorResponse = map[string]map[string]string{}
        errorResponse["hasError"] = map[string]string{"ack": "false", "permissions": "none", "error": "Not enough permissions"}
        n.Data["json"] = errorResponse
    } else {
        data, err := models.PingCluster(n.Ctx.Input.Header("user"))
        n.Data["json"] = data
        if err != nil {
            var errorResponse = map[string]map[string]string{}
            errorResponse["hasError"] = map[string]string{"ack": "false", "error": err.Error()}
            n.Data["json"] = errorResponse
        }
    }
    n.ServeJSON()
}

// @Title EditClusterValue
// @Description Edit Zeek status
// @Success 200 {object} models.zeek
// @router /editClusterValue [put]
func (n *ZeekController) EditClusterValue() {
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"EditClusterValue"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /editClusterValue [put]"
        err := models.EditClusterValue(anode, n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"DeleteClusterValue"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "DELETE"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /deleteClusterValue [delete]"
        err := models.DeleteClusterValue(anode, n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"SyncCluster"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        anode := map[string]string{}
        // var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        anode["action"] = "PUT"
        anode["controller"] = "ZEEK"
        anode["router"] = "@router /syncCluster [put]"
        err := models.SyncCluster(anode, n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"SavePolicyFiles"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        var anode map[string]map[string][]byte
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        logs.Info("ACTION -> POST")
        logs.Info("CONTROLLER -> ZEEK")
        logs.Info("ROUTER -> @router / [post]")
        for key := range anode {
            logs.Info("key -> " + key)
        }

        err := models.SavePolicyFiles(anode, n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"SyncClusterFile"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        anode := map[string][]byte{}
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        logs.Info("ACTION -> POST")
        logs.Info("CONTROLLER -> ZEEK")
        logs.Info("ROUTER -> @router / [post]")
        for key := range anode {
            logs.Info("key -> " + key)
        }

        err := models.SyncClusterFile(anode, n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"LaunchZeekMainConf"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        anode := map[string]string{}
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        logs.Info("ACTION -> POST")
        logs.Info("CONTROLLER -> ZEEK")
        logs.Info("ROUTER -> @router /LaunchZeekMainConf [put]")
        for key := range anode {
            logs.Info("key -> " + key)
        }

        err := models.LaunchZeekMainConf(anode, n.Ctx.Input.Header("user"))
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
    errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
    if errToken != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
        n.ServeJSON()
        return
    }
    permissions := []string{"SyncZeekValues"}
    hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
    if permissionsErr != nil || hasPermission == false {
        n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
    } else {
        anode := map[string]string{}
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)

        logs.Info("ACTION -> POST")
        logs.Info("CONTROLLER -> ZEEK")
        logs.Info("ROUTER -> @router /SyncZeekValues [put]")
        for key := range anode {
            logs.Info("key -> " + key)
        }

        err := models.SyncZeekValues(anode, n.Ctx.Input.Header("user"))
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}
