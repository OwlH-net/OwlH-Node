package controllers

import (
	"owlhnode/models"
	"encoding/json"
	// "strconv"
	"github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
)

type SuricataController struct {
	beego.Controller
}

// @Title GetSuricata
// @Description get Surucata status
// @Success 200 {object} models.suricata
// @router / [get]
func (m *SuricataController) Get() {
    logs.Info ("Suricata controller -> GET")
	mstatus := models.GetSuricata()
    //m.Data["json"] = map[string]string{"status": strconv.FormatBool(mstatus)}
    m.Data["json"] = mstatus
    m.ServeJSON()
}


/*
// @Title Get Suricata BPF
// @Description get Surucata BPF
// @Success 200 {object} models.suricata
// @router /bpf [get]
func (m *SuricataController) GetBPF() {
    logs.Info ("Suricata controller -> GET BPF")
    currentBPF := models.GetBPF()
    m.Data["json"] = map[string]string{"current": currentBPF}
    m.ServeJSON()
}
*/



// @Title PUT Suricata BPF
// @Description Set Surucata BPF into filter.bpf file
// @Success 200 {object} models.suricata
// @router /bpf [put]
func (n *SuricataController) SetBPF() {
    logs.Info ("Suricata controller -> SET BPF")
    
    var anode map[string]string
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    data,err := models.SetBPF(anode)

    n.Data["json"] = map[string]string{"bpf Node": data}

    if err != nil {
        logs.Info("BPF JSON RECEIVED -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }

    //newBPF := m.Ctx.Input.Param(":bpf")
    //isSetBPF := models.SetBPF(newBPF)
    n.Data["json"] = map[string]string{"status": "true"}
    n.ServeJSON()
}

// @Title RetrieveFile
// @Description Retrieve file from master
// @Success 200 {object} models.Node
// @Failure 403 body is empty
// @router /retrieve [put]
func (n *SuricataController) RetrieveFile() {
    logs.Info("retrieve -> In")
    var anode map[string][]byte
    json.Unmarshal(n.Ctx.Input.RequestBody, &anode)
    err := models.RetrieveFile(anode)
    //logs.Info(string(anode["data"]))
    n.Data["json"] = map[string]string{"ack": "true"}
    if err != nil {
        logs.Info("Ruleset retrieve OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("retrieve -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title RunSuricata
// @Description Run suricata system
// @Success 200 {object} models.suricata
// @Failure 403 body is empty
// @router /RunSuricata [put]
func (n *SuricataController) RunSuricata() {
    logs.Info("RunSuricata -> In")
    data,err := models.RunSuricata()
    n.Data["json"] = data
    if err != nil {
        logs.Info("RunSuricata OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("RunSuricata -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}

// @Title StopSuricata
// @Description Run suricata system
// @Success 200 {object} models.suricata
// @Failure 403 body is empty
// @router /StopSuricata [put]
func (n *SuricataController) StopSuricata() {
    logs.Info("StopSuricata -> In")
    data,err := models.StopSuricata()
    n.Data["json"] = data
    if err != nil {
        logs.Info("StopSuricata OUT -- ERROR : %s", err.Error())
        n.Data["json"] = map[string]string{"ack": "false", "error": err.Error()}
    }
    logs.Info("StopSuricata -> OUT -> %s", n.Data["json"])
    n.ServeJSON()
}