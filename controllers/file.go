package controllers

import (
    "owlhnode/models"
    "encoding/json"
    "github.com/astaxie/beego"
    "owlhnode/validation"
    "github.com/astaxie/beego/logs"
)

type FileController struct {
    beego.Controller
}

// @Title SendFile
// @Description send back the requested file from master for show on webpage "edit.html"
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /:fileName [get]
func (n *FileController) SendFile() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("File Error validating token from master")
logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        fileName := n.GetString(":fileName")
        data, err := models.SendFile(fileName)
    
        n.Data["json"] = data
        if err != nil {
            logs.Info("send OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        logs.Info("send -> OUT -> %s", n.Data["json"])
    }
    n.ServeJSON()
}

// @Title SaveFile
// @Description save changes over requested file on webpage "edit.html"
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router / [put]
func (n *FileController) SaveFile() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "put")
    if err != nil {
        logs.Error("File Error validating token from master")
logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        var anode map[string]string
        json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
        err := models.SaveFile(anode)
        anode["action"] = "PUT"
        anode["controller"] = "FILE"
        anode["router"] = "@router / [put]"
    
        n.Data["json"] = map[string]string{"ack": "true"}
        if err != nil {
            logs.Info("save OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        logs.Info("save -> OUT -> %s", n.Data["json"])
    }
    n.ServeJSON()
}

// @Title GetAllFiles
// @Description save changes over requested file on webpage "edit.html"
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router / [get]
func (n *FileController) GetAllFiles() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("File Error validating token from master")
logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        data,err := models.GetAllFiles()
    
        n.Data["json"] = data
        if err != nil {
            logs.Info("GetAllFiles OUT -- ERROR : %s", err.Error())
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
        logs.Info("GetAllFiles return %s", n.Data["json"])
    }
    n.ServeJSON()
}

// @Title ReloadFilesData
// @Description load new files size
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /reloadFilesData [get]
func (n *FileController) ReloadFilesData() {
    permissions,err := validation.CheckToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"), n.Ctx.Input.Header("uuid"), "get")
    if err != nil {
        logs.Error("File Error validating token from master")
logs.Error(err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "token":"none"}
    }else if !permissions{
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error(), "permissions":"none"}
    }else{         
        data,err := models.ReloadFilesData()
    
        n.Data["json"] = data
        if err != nil {
            n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
        }
    }
    n.ServeJSON()
}