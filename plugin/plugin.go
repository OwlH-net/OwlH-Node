package plugin

import (
    "github.com/astaxie/beego/logs"	
	"owlhnode/database"
    // "owlhnode/suricata"
    "os/exec"
    "bytes"
    "os"
    "strconv"
    "io/ioutil"
	"owlhnode/utils"
)

func ChangeServiceStatus(anode map[string]string)(err error) {  
    if anode["status"] == "enabled"{
        iface,err := ndb.GetPluginsByParam(anode["server"], "interface")
        if err != nil {logs.Error("plugin/ChangeServiceStatus error GetPluginsByParam: "+err.Error()); return err}

        cmd := exec.Command("suricata", "-D", "-c", "/etc/suricata/suricata.yaml", "-i", iface, "--pidfile", "/etc/suricata/pidfile/"+anode["server"]+"-pidfile.txt")
        var stdBuffer bytes.Buffer
        cmd.Stderr = &stdBuffer

        err = cmd.Run()
        if err != nil {
            logs.Error(stdBuffer.String())
            logs.Error("plugin/ChangeServiceStatus error launching Suricata: "+err.Error()); 
            //delete pid file
            err = os.Remove("/etc/suricata/pidfile/"+anode["server"]+"-pidfile.txt")
            if err != nil {logs.Error("plugin/SaveSuricataInterface error deleting a pid file: "+err.Error())}
        }else{
            //read file
            currentpid, err := os.Open("/etc/suricata/pidfile/"+anode["server"]+"-pidfile.txt")
            if err != nil {logs.Error("plugin/ChangeServiceStatus error launching Suricata: "+err.Error()); return err}
            defer currentpid.Close()
            pid, err := ioutil.ReadAll(currentpid)

            //save pid to db
            err = ndb.UpdatePluginValue(anode["server"],"pid",string(pid))
            if err != nil {logs.Error("plugin/SaveSuricataInterface error updating pid at DB: "+err.Error()); return err}
            
            //change DB status
            err = ndb.UpdatePluginValue(anode["server"],anode["param"],anode["status"])
            if err != nil {logs.Error("plugin/ChangeServiceStatus error: "+err.Error()); return err}
        }
    }else if anode["status"] == "disabled"{
        //pid
        currentpid, err := os.Open("/etc/suricata/pidfile/"+anode["server"]+"-pidfile.txt")
        if err != nil {logs.Error("plugin/ChangeServiceStatus error reading Suricata pid: "+err.Error()); return err}
        defer currentpid.Close()
        pid, err := ioutil.ReadAll(currentpid)
        PidInt,_ := strconv.Atoi(string(pid))

        process, err := os.FindProcess(PidInt)
        err = process.Kill()
        if err != nil {logs.Error("plugin/ChangeServiceStatus error killing Suricata process: "+err.Error()); return err}

        // //kill
        // killPID := exec.Command("bash","-c","kill","-9", string(pid))
        // err = killPID.Run()
        // if err != nil {logs.Error("plugin/ChangeServiceStatus error stopping Suricata: "+err.Error()); return err}

        //delete pid file
        err = os.Remove("/etc/suricata/pidfile/"+anode["server"]+"-pidfile.txt")
        if err != nil {logs.Error("plugin/SaveSuricataInterface error deleting a pid file: "+err.Error()); return err}

        //change DB pid
        err = ndb.UpdatePluginValue(anode["service"],"pid","none") 
        if err != nil {logs.Error("plugin/SaveSuricataInterface error updating pid at DB: "+err.Error()); return err}

        //change DB status
        err = ndb.UpdatePluginValue(anode["server"],anode["param"],anode["status"])
        if err != nil {logs.Error("plugin/ChangeServiceStatus error: "+err.Error()); return err}
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
        err = ndb.InsertPluginService(uuid, "pid", "none"); if err != nil {logs.Error("InsertPluginService ruleset Error: "+err.Error()); return err}    
    }

    return nil
}

func SaveSuricataInterface(anode map[string]string)(err error) {
    err = ndb.UpdatePluginValue(anode["service"],anode["param"],anode["interface"])
    if err != nil {logs.Error("plugin/SaveSuricataInterface error: "+err.Error()); return err}
    return err
}