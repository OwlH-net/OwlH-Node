package controllers

import (
	"encoding/json"

	"github.com/OwlH-net/OwlH-Node/models"
	"github.com/OwlH-net/OwlH-Node/validation"
	"github.com/astaxie/beego"
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
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"SendFile"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		fileName := n.GetString(":fileName")
		data, err := models.SendFile(fileName, n.Ctx.Input.Header("user"))
		n.Data["json"] = data
		if err != nil {
			logs.Error("send OUT -- ERROR : %s", err.Error())
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
	}
	n.ServeJSON()
}

// @Title SaveFile
// @Description save changes over requested file on webpage "edit.html"
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router / [put]
func (n *FileController) SaveFile() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"SaveFile"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		var anode map[string]string
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		err := models.SaveFile(anode, n.Ctx.Input.Header("user"))
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
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"GetAllFiles"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		data, err := models.GetAllFiles(n.Ctx.Input.Header("user"))

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
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"StopStapServiceNode"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		data, err := models.ReloadFilesData(n.Ctx.Input.Header("user"))

		n.Data["json"] = data
		if err != nil {
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
	}
	n.ServeJSON()
}
