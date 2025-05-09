package controllers

import (
	"encoding/json"

	"github.com/OwlH-net/OwlH-Node/models"
	"github.com/OwlH-net/OwlH-Node/validation"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type WazuhController struct {
	beego.Controller
}

type WazuhDetails struct {
	Path    bool   `json:"path"`
	Bin     bool   `json:"bin"`
	Running bool   `json:"running"`
	Name    string `json:"name"`
	ID      string `json:"id"`
	Ip      string `json:"ip"`
}

// @Title GetWazuh
// @Description get Wazuh status
// @Success 200 {object} models.wazuh
// @router / [get]
func (n *WazuhController) Get() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"GetWazuh"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		mstatus, err := models.GetWazuh(n.Ctx.Input.Header("user"))
		if err == nil {
			var wazuhDetails WazuhDetails
			wazuhDetails.Bin = mstatus["bin"]
			wazuhDetails.Path = mstatus["path"]
			wazuhDetails.Running = mstatus["running"]
			details, err := models.GetWazuhDetails()
			if err == nil {
				wazuhDetails.ID = details.ID
				wazuhDetails.Ip = details.Ip
				wazuhDetails.Name = details.Name
			}
			logs.Info("wazuh -> to convert -> %+v", wazuhDetails)

			returnData, err := json.Marshal(wazuhDetails)

			n.Data["json"] = string(returnData)
			logs.Info("wazuh -> returning -> %+v", string(returnData))

		} else {
			logs.Info("GetWazuh OUT -- ERROR : %s", err.Error())
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
	}
	n.ServeJSON()
}

// @Title RunWazuh
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /RunWazuh [put]
func (n *WazuhController) RunWazuh() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"RunWazuh"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		logs.Info("RunWazuh -> In")
		var anode = map[string]string{}
		anode["action"] = "PUT"
		anode["controller"] = "WAZUH"
		anode["router"] = "@router /RunWazuh [put]"
		logs.Info("============")
		logs.Info("WAZUH - RunWazuh")
		for key := range anode {
			logs.Info(key + " -> " + anode[key])
		}
		data, err := models.RunWazuh(n.Ctx.Input.Header("user"))
		n.Data["json"] = data
		if err != nil {
			logs.Info("RunWazuh OUT -- ERROR : %s", err.Error())
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
		logs.Info("RunWazuh -> OUT -> %s", n.Data["json"])
	}
	n.ServeJSON()
}

// @Title StopWazuh
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /StopWazuh [put]
func (n *WazuhController) StopWazuh() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"StopWazuh"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		logs.Info("StopWazuh -> In")
		var anode = map[string]string{}
		anode["action"] = "PUT"
		anode["controller"] = "WAZUH"
		anode["router"] = "@router /StopWazuh [put]"
		logs.Info("============")
		logs.Info("WAZUH - StopWazuh")
		for key := range anode {
			logs.Info(key + " -> " + anode[key])
		}
		data, err := models.StopWazuh(n.Ctx.Input.Header("user"))
		n.Data["json"] = data
		if err != nil {
			logs.Info("StopWazuh OUT -- ERROR : %s", err.Error())
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
		logs.Info("StopWazuh -> OUT -> %s", n.Data["json"])
	}
	n.ServeJSON()
}

// @Title PingWazuhFiles
// @Description get Wazuh status
// @Success 200 {object} models.wazuh
// @router /pingWazuhFiles [get]
func (n *WazuhController) PingWazuhFiles() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		var errorResponse = map[string]map[string]string{}
		errorResponse["hasError"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.Data["json"] = errorResponse
		n.ServeJSON()
		return
	}
	permissions := []string{"PingWazuhFiles"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		var errorResponse = map[string]map[string]string{}
		errorResponse["hasError"] = map[string]string{"ack": "false", "permissions": "none", "error": "Not enough permissions"}
		n.Data["json"] = errorResponse
	} else {
		files, err := models.PingWazuhFiles(n.Ctx.Input.Header("user"))
		n.Data["json"] = files
		if err != nil {
			var errorResponse = map[string]map[string]string{}
			errorResponse["hasError"] = map[string]string{"ack": "false", "error": err.Error()}
			n.Data["json"] = errorResponse
		}
	}
	n.ServeJSON()
}

// @Title DeleteWazuhFile
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /deleteWazuhFile [delete]
func (n *WazuhController) DeleteWazuhFile() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"DeleteWazuhFile"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		var anode map[string]interface{}
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		anode["action"] = "PUT"
		anode["controller"] = "SURICATA"
		anode["router"] = "@router /StopSuricata [put]"
		err := models.DeleteWazuhFile(anode, n.Ctx.Input.Header("user"))
		n.Data["json"] = map[string]string{"ack": "true"}

		if err != nil {
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
	}
	n.ServeJSON()
}

// @Title AddWazuhFile
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /addWazuhFile [put]
func (n *WazuhController) AddWazuhFile() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"AddWazuhFile"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		var anode map[string]interface{}
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		anode["action"] = "PUT"
		anode["controller"] = "SURICATA"
		anode["router"] = "@router /StopSuricata [put]"
		err := models.AddWazuhFile(anode, n.Ctx.Input.Header("user"))
		n.Data["json"] = map[string]string{"ack": "true"}

		if err != nil {
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
	}
	n.ServeJSON()
}

// @Title LoadFileLastLines
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /loadFileLastLines [put]
func (n *WazuhController) LoadFileLastLines() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"LoadFileLastLines"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		var anode = map[string]string{}
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		anode["action"] = "PUT"
		anode["controller"] = "SURICATA"
		anode["router"] = "@router /StopSuricata [put]"
		data, err := models.LoadFileLastLines(anode, n.Ctx.Input.Header("user"))
		n.Data["json"] = data

		if err != nil {
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
	}
	n.ServeJSON()
}

// @Title SaveFileContentWazuh
// @Description Run wazuh system
// @Success 200 {object} models.wazuh
// @Failure 403 body is empty
// @router /saveFileContentWazuh [put]
func (n *WazuhController) SaveFileContentWazuh() {
	errToken := validation.VerifyToken(n.Ctx.Input.Header("token"), n.Ctx.Input.Header("user"))
	if errToken != nil {
		n.Data["json"] = map[string]string{"ack": "false", "error": errToken.Error(), "token": "none"}
		n.ServeJSON()
		return
	}
	permissions := []string{"SaveFileContentWazuh"}
	hasPermission, permissionsErr := validation.VerifyPermissions(n.Ctx.Input.Header("user"), "any", permissions)
	if permissionsErr != nil || hasPermission == false {
		n.Data["json"] = map[string]string{"ack": "false", "permissions": "none"}
	} else {
		var anode = map[string]string{}
		json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
		anode["action"] = "PUT"
		anode["controller"] = "SURICATA"
		anode["router"] = "@router /StopSuricata [put]"

		err := models.SaveFileContentWazuh(anode, n.Ctx.Input.Header("user"))
		n.Data["json"] = map[string]string{"ack": "true"}

		if err != nil {
			n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
		}
	}
	n.ServeJSON()
}
