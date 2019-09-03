package plugin

import (
    "github.com/astaxie/beego/logs"	
	"owlhnode/database"
    // "owlhnode/suricata"
    "os/exec"
	"owlhnode/utils"
)

func ChangeServiceStatus(anode map[string]string)(err error) {
	err = ndb.UpdatePluginValue(anode["server"],anode["param"],anode["status"])
    if err != nil {logs.Error("plugin/ChangeServiceStatus error: "+err.Error()); return err}
    
    if anode["status"] == "enabled"{
        iface,err := ndb.GetPluginsByParam(anode["server"], "interface")
        if err != nil {logs.Error("plugin/ChangeServiceStatus error GetPluginsByParam: "+err.Error()); return err}

        cmd := "suricata -c /etc/suricata/suricata.yaml -i "+iface+" &"
        output, err := exec.Command("bash", "-c", cmd).Output()
        if err != nil {logs.Error("plugin/ChangeServiceStatus error launching Suricata: "+err.Error()); return err}
        logs.Debug(string(output))
        logs.Debug(string(output))
        logs.Debug(string(output))
        logs.Debug(string(output))
        logs.Debug(string(output))
        logs.Debug(string(output))
        
        logs.Debug(string(output))
    }else if anode["status"] == "enabled"{
        //kill -9 pidof 65465456
    }
    

    return err
}

func ChangeMainServiceStatus(anode map[string]string)(err error) {
	err = ndb.UpdateMainconfValue(anode["service"],anode["param"],anode["status"])
	if err != nil {logs.Error("plugin/ChangeMainServiceStatus error: "+err.Error()); return err}
    return err
}

func DeleteService(anode map[string]string)(err error) {
	err = ndb.DeleteService(anode["server"])
	if err != nil {logs.Error("plugin/DeleteService error: "+err.Error()); return err}
    return err
}

func AddPluginService(anode map[string]string) (err error) {
    uuid := utils.Generate()
    err = ndb.InsertPluginService(uuid, "status", "disabled"); if err != nil {logs.Error("InsertPluginService status Error: "+err.Error()); return err}
    err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
    err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
    err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}    
    if anode["type"] == "suricata"{
        err = ndb.InsertPluginService(uuid, "interface", "none"); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "bpf", "none"); if err != nil {logs.Error("InsertPluginService bpf Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "ruleset", "none"); if err != nil {logs.Error("InsertPluginService ruleset Error: "+err.Error()); return err}    
    }

    return nil
}

func SaveSuricataInterface(anode map[string]string)(err error) {
    err = ndb.UpdatePluginValue(anode["service"],anode["param"],anode["interface"])
    if err != nil {logs.Error("plugin/SaveSuricataInterface error: "+err.Error()); return err}
    return err
}