package controllers

import (
	"owlhnode/models"
	"encoding/json"
	"strconv"
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
	m.Data["json"] = map[string]string{"status": strconv.FormatBool(mstatus)}
    m.ServeJSON()
}

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