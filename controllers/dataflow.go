package controllers

import (
	"owlhnode/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type DataflowController struct {
	beego.Controller
}

// @Title ChangeDataflowValues
// @Description Make a deploy for selected action
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /changeDataflowValues [put]
func (n *DataflowController) ChangeDataflowValues() {
	var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)	
    err := models.ChangeDataflowValues(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title LoadDataflowValues
// @Description Load data flow values for
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /loadDataflowValues [get]
func (n *DataflowController) LoadDataflowValues() {
    data,err := models.LoadDataflowValues()
    n.Data["json"] = data
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}