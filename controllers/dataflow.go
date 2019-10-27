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
	anode["action"] = "PUT"
    anode["controller"] = "DATAFLOW"
    anode["router"] = "@router /changeDataflowValues [put]"

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

// @Title SaveSocketToNetwork
// @Description Save socket information to Network at node dataflow
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /saveSocketToNetwork [put]
func (n *DataflowController) SaveSocketToNetwork() {
	var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)	
	anode["action"] = "PUT"
    anode["controller"] = "DATAFLOW"
    anode["router"] = "@router /saveSocketToNetwork [put]"    
    err := models.SaveSocketToNetwork(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title SaveNewLocal
// @Description Save socket information to Network at node dataflow
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /saveNewLocal [put]
func (n *DataflowController) SaveNewLocal() {
	var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)	
	anode["action"] = "PUT"
    anode["controller"] = "DATAFLOW"
    anode["router"] = "@router /saveNewLocal [put]"    
    err := models.SaveNewLocal(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title SaveVxLAN
// @Description Save socket information to Network at node dataflow
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /saveVxLAN [put]
func (n *DataflowController) SaveVxLAN() {
	var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)	
	anode["action"] = "PUT"
    anode["controller"] = "DATAFLOW"
    anode["router"] = "@router /saveVxLAN [put]"    
    err := models.SaveVxLAN(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title SaveSocketToNetworkSelected
// @Description Save socket information to Network at node dataflow
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /saveSocketToNetworkSelected [put]
func (n *DataflowController) SaveSocketToNetworkSelected() {
	var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)	
	anode["action"] = "PUT"
    anode["controller"] = "DATAFLOW"
    anode["router"] = "@router /saveSocketToNetworkSelected [put]"    
    err := models.SaveSocketToNetworkSelected(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}

// @Title DeleteDataFlowValueSelected
// @Description Save socket information to Network at node dataflow
// @Success 200 {object} models.file
// @Failure 403 body is empty
// @router /deleteDataFlowValueSelected [delete]
func (n *DataflowController) DeleteDataFlowValueSelected() {
	var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)	
	anode["action"] = "DELETE"
    anode["controller"] = "DATAFLOW"
    anode["router"] = "@router /deleteDataFlowValueSelected [delete]"    
    err := models.DeleteDataFlowValueSelected(anode)
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    n.ServeJSON()
}