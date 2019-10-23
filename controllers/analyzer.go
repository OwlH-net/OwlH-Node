package controllers

import (
	"owlhnode/models"
	"owlhnode/analyzer"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type AnalyzerController struct {
	beego.Controller
}

// @Title PingAnalyzer
// @Description PingAnalyzer status
// @Success 200 {object} models.analyzer
// @router /pingAnalyzer [get]
func (m *AnalyzerController) PingAnalyzer() {	
	data, err := models.PingAnalyzer()
	m.Data["json"] = data
	if err != nil {
        logs.Info("PingAnalyzer OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}
    m.ServeJSON()
}

// @Title ChangeAnalyzerStatus
// @Description ChangeAnalyzerStatus status
// @Success 200 {object} models.analyzer
// @router /changeAnalyzerStatus [put]
func (m *AnalyzerController) ChangeAnalyzerStatus() {	
	var anode map[string]string
    json.Unmarshal(m.Ctx.Input.RequestBody, &anode)
	anode["action"] = "PUT"
    anode["controller"] = "ANALYZER"
    anode["router"] = "@router /changeAnalyzerStatus [put]"
    err := models.ChangeAnalyzerStatus(anode)

	m.Data["json"] = map[string]string{"ack": "true"}
	if err != nil {
        logs.Info("ChangeAnalyzerStatus OUT -- ERROR : %s", err.Error())
        m.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
	}else{
		analyzer.Init()
	}
    m.ServeJSON()
}