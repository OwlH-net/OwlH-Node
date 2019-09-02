package plugin

import (
    "github.com/astaxie/beego/logs"	
	"owlhnode/database"
	"owlhnode/suricata"
	"owlhnode/utils"
)

func ChangeServiceStatus(anode map[string]string)(err error) {
	err = ndb.UpdatePluginValue(anode["server"],anode["param"],anode["status"])
    if err != nil {logs.Error("plugin/ChangeServiceStatus error: "+err.Error()); return err}
    
    if anode["status"] == "enabled"{
        status,err :=suricata.RunSuricata()
        if err != nil {logs.Error("plugin/RunSuricata error: "+err.Error()); return err}    
        logs.Debug(status)
        logs.Debug(status)
        logs.Debug(status)
    }else if anode["status"] == "enabled"{
        status,err :=suricata.StopSuricata()
        if err != nil {logs.Error("plugin/StopSuricata error: "+err.Error()); return err}    
        logs.Debug(status)
        logs.Debug(status)
        logs.Debug(status)
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
    err = ndb.InsertPluginService(uuid, "interface", "none"); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
    err = ndb.InsertPluginService(uuid, "bpf", "none"); if err != nil {logs.Error("InsertPluginService bpf Error: "+err.Error()); return err}
    err = ndb.InsertPluginService(uuid, "ruleset", "none"); if err != nil {logs.Error("InsertPluginService ruleset Error: "+err.Error()); return err}    
    err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}    

    return nil
}