package plugin

import (
    "github.com/astaxie/beego/logs"	
	"owlhnode/database"
    // "owlhnode/suricata"
    "os/exec"
    "bytes"
    "errors"
    "os"
    "strconv"
    "strings"
    "io/ioutil"
	"owlhnode/utils"
)

func ChangeServiceStatus(anode map[string]string)(err error) {  
    if anode["type"] == "suricata"{
        if anode["status"] == "enabled"{
            allPlugins,err := ndb.GetPlugins()
            for x := range allPlugins {
                //get all db values and check if any
                if allPlugins[x]["pid"] != "none" && allPlugins[x]["interface"] == anode["interface"] && allPlugins[x]["status"] == "enabled" && x != anode["server"]{
                    logs.Error("Can't launch more than one suricata with same interface. Please, select other interface.")
                    return errors.New("Can't launch more than one suricata with same interface. Please, select other interface.")
                }
            }
            err = LaunchSuricataService(anode["server"], anode["interface"])
            if err != nil {logs.Error("LaunchSuricataService status Error: "+err.Error()); return err}
        }else if anode["status"] == "disabled"{
            err = StopSuricataService(anode["server"], anode["status"])
            if err != nil {logs.Error("StopSuricataService status Error: "+err.Error()); return err}

        }
    } else if anode["type"] == "zeek"{
        mainConfData, err := ndb.GetMainconfData()
        if (mainConfData["zeek"]["status"] == "disabled"){ return nil }

        allPlugins,err := ndb.GetPlugins()
        for x := range allPlugins {
            if anode["status"] == "enabled" && allPlugins[x]["type"] == "zeek"{
                err = ndb.UpdatePluginValue(x,"previousStatus","none")
                if err != nil {logs.Error("plugin/ChangeServiceStatus error updating pid at DB: "+err.Error()); return err}
                
                err = ndb.UpdatePluginValue(x,"status","enabled")
                if err != nil {logs.Error("plugin/ChangeServiceStatus error updating pid at DB: "+err.Error()); return err}
            }else if anode["status"] == "disabled" && allPlugins[x]["type"] == "zeek"{
                err = ndb.UpdatePluginValue(x,"previousStatus",allPlugins[x]["status"])
                if err != nil {logs.Error("plugin/ChangeServiceStatus error updating pid at DB: "+err.Error()); return err}
                
                err = ndb.UpdatePluginValue(x,"status","disabled")
                if err != nil {logs.Error("plugin/ChangeServiceStatus error updating pid at DB: "+err.Error()); return err}
            }
        }
        
    }
    return err
}

func ChangeMainServiceStatus(anode map[string]string)(err error) {    
    err = ndb.UpdateMainconfValue(anode["service"],anode["param"],anode["status"])
    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error: "+err.Error()); return err}

    allPlugins,err := ndb.GetPlugins()
    if anode["service"] == "suricata" {
        for x := range allPlugins {
            if anode["status"] == "disabled"{
                if allPlugins[x]["status"] == "enabled" && allPlugins[x]["type"] == "suricata"{
                    err = StopSuricataService(x, allPlugins[x]["status"])
                    if err != nil {logs.Error("StopSuricataService status Error: "+err.Error()); return err}
                }
            }else if anode["status"] == "enabled"{
                if allPlugins[x]["previousStatus"] == "enabled" && allPlugins[x]["type"] == "suricata"{
                    err = LaunchSuricataService(x, allPlugins[x]["interface"])
                    if err != nil {logs.Error("LaunchSuricataService status Error: "+err.Error()); return err}
                }
            }
        }
    } else if anode["service"] == "zeek" {
        for x := range allPlugins {
            if anode["status"] == "disabled"{
                if allPlugins[x]["status"] == "enabled" && allPlugins[x]["type"] == "zeek"{
                    err = ndb.UpdatePluginValue(x,"previousStatus",allPlugins[x]["status"])
                    if err != nil {logs.Error("plugin/SaveSuricataInterface error updating pid at DB: "+err.Error()); return err}
                    
                    err = ndb.UpdatePluginValue(x,"status","disabled")
                    if err != nil {logs.Error("plugin/SaveSuricataInterface error updating pid at DB: "+err.Error()); return err}
                }
            }else if anode["status"] == "enabled"{
                if allPlugins[x]["previousStatus"] == "enabled" && allPlugins[x]["type"] == "zeek"{
                    err = ndb.UpdatePluginValue(x,"previousStatus","none")
                    if err != nil {logs.Error("plugin/SaveSuricataInterface error updating pid at DB: "+err.Error()); return err}
                    
                    err = ndb.UpdatePluginValue(x,"status",allPlugins[x]["previousStatus"])
                    if err != nil {logs.Error("plugin/SaveSuricataInterface error updating pid at DB: "+err.Error()); return err}
                }
            }
        }
    }

    return err
}

func DeleteService(anode map[string]string)(err error) {
	err = ndb.DeleteService(anode["server"])
    if err != nil {logs.Error("plugin/DeleteService error: "+err.Error()); return err}

    if _, err := os.Stat("/etc/suricata/bpf/"+anode["server"]+" - filter.bpf"); !os.IsNotExist(err) {
        err = os.Remove("/etc/suricata/bpf/"+anode["server"]+" - filter.bpf")
        if err != nil {logs.Error("plugin/SaveSuricataInterface error deleting a pid file: "+err.Error())}
    }

    return err
}

func AddPluginService(anode map[string]string) (err error) {
    uuid := utils.Generate()
    err = ndb.InsertPluginService(uuid, "status", "disabled"); if err != nil {logs.Error("InsertPluginService status Error: "+err.Error()); return err}
    err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
    err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
    err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}    
    err = ndb.InsertPluginService(uuid, "previousStatus", "none"); if err != nil {logs.Error("InsertPluginService previousStatus Error: "+err.Error()); return err}    
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

func LaunchSuricataService(uuid string, iface string)(err error){

    mainConfData, err := ndb.GetMainconfData()
    if (mainConfData["suricata"]["status"] == "disabled"){ return nil }

    cmd := exec.Command("suricata", "-D", "-c", "/etc/suricata/suricata.yaml", "-i", iface, "--pidfile", "/var/run/suricata/"+uuid+"-pidfile.pid")
    var stdBuffer bytes.Buffer
    cmd.Stderr = &stdBuffer

    err = cmd.Run()
    if err != nil {
        logs.Error(stdBuffer.String())
        logs.Error("plugin/LaunchSuricataService error launching Suricata: "+err.Error()); 
        //delete pid file
        err = os.Remove("/var/run/suricata/"+uuid+"-pidfile.pid")
        if err != nil {logs.Error("plugin/SaveSuricataInterface error deleting a pid file: "+err.Error())}
    }else{
        //read file
        currentpid, err := os.Open("/var/run/suricata/"+uuid+"-pidfile.pid")
        if err != nil {logs.Error("plugin/LaunchSuricataService error openning Suricata: "+err.Error()); return err}
        defer currentpid.Close()
        pid, err := ioutil.ReadAll(currentpid)

        //save pid to db
        err = ndb.UpdatePluginValue(uuid,"pid",string(pid))
        if err != nil {logs.Error("plugin/SaveSuricataInterface error updating pid at DB: "+err.Error()); return err}
        
        //change DB status
        err = ndb.UpdatePluginValue(uuid,"previousStatus","none")
        if err != nil {logs.Error("plugin/LaunchSuricataService error: "+err.Error()); return err}

        //change DB status
        err = ndb.UpdatePluginValue(uuid,"status","enabled")
        if err != nil {logs.Error("plugin/LaunchSuricataService error: "+err.Error()); return err}
    }
    return nil
}

func StopSuricataService(uuid string, status string)(err error){
    //pid
    // currentpid, err := os.Open("/var/run/suricata/"+uuid+"-pidfile.pid")
    // if err != nil {logs.Error("plugin/ChangeServiceStatus error reading Suricata pid: "+err.Error()); return err}
    // defer currentpid.Close()
    // pid, err := ioutil.ReadAll(currentpid)
    allPlugin,err := ndb.GetPlugins()
    
    //kill suricata process
    PidInt,_ := strconv.Atoi(strings.Trim(string(allPlugin[uuid]["pid"]), "\n"))
    process, _ := os.FindProcess(PidInt)
    _ = process.Kill()
    // if err != nil {logs.Error("plugin/StopSuricataService error killing Suricata process: "+err.Error()); return err}
    
    //delete pid file
    _ = os.Remove("/var/run/suricata/"+uuid+"-pidfile.pid")
    // if err != nil {logs.Error("plugin/SaveSuricataInterface error deleting a pid file: "+err.Error()); return err}

    //change DB pid
    err = ndb.UpdatePluginValue(uuid,"pid","none") 
    if err != nil {logs.Error("plugin/SaveSuricataInterface error updating pid at DB: "+err.Error()); return err}

    //change DB status
    err = ndb.UpdatePluginValue(uuid,"previousStatus",status)
    if err != nil {logs.Error("plugin/StopSuricataService error: "+err.Error()); return err}

    //change DB status
    err = ndb.UpdatePluginValue(uuid,"status","disabled")
    if err != nil {logs.Error("plugin/StopSuricataService error: "+err.Error()); return err}

    return nil
}