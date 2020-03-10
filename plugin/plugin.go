package plugin

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
    "owlhnode/zeek"
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
    allPlugins,err := ndb.GetPlugins()
    if anode["type"] == "suricata"{
        if anode["status"] == "enabled"{
            for x := range allPlugins {
                //get all db values and check if any
                if allPlugins[x]["type"] == "suricata" && allPlugins[x]["pid"] != "none" && allPlugins[x]["interface"] == anode["interface"] && allPlugins[x]["status"] == "enabled" && x != anode["service"]{
                    logs.Error("Can't launch more than one suricata with same interface. Please, select other interface.")
                    return errors.New("Can't launch more than one suricata with same interface. Please, select other interface.")
                }
            }
            err = LaunchSuricataService(anode["service"], anode["interface"])
            if err != nil {logs.Error("LaunchSuricataService status Error: "+err.Error()); return err}
        }else if anode["status"] == "disabled"{
            err = StopSuricataService(anode["service"], anode["status"])
            if err != nil {logs.Error("StopSuricataService status Error: "+err.Error()); return err}

        }
    } else if anode["type"] == "zeek"{
        mainConfData, err := ndb.GetMainconfData()
        if (mainConfData["zeek"]["status"] == "disabled"){ return nil }        
        if anode["status"] == "enabled"{
            err = zeek.DeployZeek()
            if err != nil {logs.Error("plugin/ChangeServiceStatus error deploying zeek: "+err.Error()); return err}

            err = ndb.UpdatePluginValue(anode["service"],"previousStatus","none")
            if err != nil {logs.Error("plugin/ChangeServiceStatus error updating zeek previousStatus to none: "+err.Error()); return err}

            err = ndb.UpdatePluginValue(anode["service"],"status","enabled")
            if err != nil {logs.Error("plugin/ChangeServiceStatus error updating zeek status to enabled: "+err.Error()); return err}
        } else if anode["status"] == "disabled"{
            data, err := zeek.StopZeek(); logs.Error(data)
            if err != nil {logs.Error("plugin/ChangeServiceStatus error deploying zeek: "+err.Error()); return err}

            err = ndb.UpdatePluginValue(anode["service"],"previousStatus",anode["status"])
            if err != nil {logs.Error("plugin/ChangeServiceStatus error updating zeek previousStatus to status: "+err.Error()); return err}

            err = ndb.UpdatePluginValue(anode["service"],"status","disabled")
            if err != nil {logs.Error("plugin/ChangeServiceStatus error updating zeek status to disabled: "+err.Error()); return err}
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
                    err = ndb.UpdatePluginValue(x,"previousStatus","enabled")
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error updating pid at DB: "+err.Error()); return err}

                    err = ndb.UpdatePluginValue(x,"status","disabled")
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error updating pid at DB: "+err.Error()); return err}
                    
                    data, err := zeek.StopZeek(); logs.Error(data)
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error deploying zeek: "+err.Error()); return err}
                }
            }else if anode["status"] == "enabled"{
                if allPlugins[x]["previousStatus"] == "enabled" && allPlugins[x]["type"] == "zeek"{
                    err = ndb.UpdatePluginValue(x,"previousStatus","none")
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error updating pid at DB: "+err.Error()); return err}

                    err = ndb.UpdatePluginValue(x,"status","enabled")
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error updating pid at DB: "+err.Error()); return err}

                    err = zeek.DeployZeek()
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error deploying zeek: "+err.Error()); return err}
                }
            }
        }
    }

    return err
}

func DeleteService(anode map[string]string)(err error) {
    path, err := utils.GetKeyValueString("suricataBPF", "pathBPF")
    if err != nil {logs.Error("DeleteService Error getting data from main.conf: "+err.Error())}
    filter, err := utils.GetKeyValueString("suricataBPF", "fileBPF")
    if err != nil {logs.Error("DeleteService Error getting data from main.conf: "+err.Error())}
    
    allPlugins,err := ndb.GetPlugins()
    if allPlugins[anode["service"]]["type"] == "suricata" {
        if allPlugins[anode["service"]]["status"] == "enabled" {
            err = StopSuricataService(anode["service"], allPlugins[anode["service"]]["status"])
            if err != nil {logs.Error("plugin/DeleteService error stopping suricata: "+err.Error()); return err}
            logs.Error("suricata 3")
        }
        if _, err := os.Stat(path+anode["service"]+"-"+filter); !os.IsNotExist(err) {
            err = os.Remove(path+anode["service"]+"-"+filter)
            if err != nil {logs.Error("plugin/SaveSuricataInterface error deleting a pid file: "+err.Error())}
        }
    }else if allPlugins[anode["service"]]["type"] == "zeek" {
        if allPlugins[anode["service"]]["status"] == "enabled" {
            _, err := zeek.StopZeek();
            if err != nil {logs.Error("plugin/DeleteService error stopping Zeek: "+err.Error())}
        }
    }else if allPlugins[anode["service"]]["type"] == "socket-network" || allPlugins[anode["service"]]["type"] == "socket-pcap" || allPlugins[anode["service"]]["type"] == "network-socket"{        
        if allPlugins[anode["service"]]["pid"] != "none" {
            err := StopStapService(anode);
            if err != nil {logs.Error("plugin/DeleteService error stopping STAP service: "+err.Error())}
        }
    }

    err = ndb.DeleteService(anode["service"])
    if err != nil {logs.Error("plugin/DeleteService error: "+err.Error()); return err}    

    return err
}

func AddPluginService(anode map[string]string) (err error) {
    uuid := utils.Generate()    
    if anode["type"] == "socket-network" || anode["type"] == "socket-pcap" || anode["type"] == "network-socket" {
        if _, err := os.Stat(anode["cert"]); os.IsNotExist(err) {
            return errors.New("STAP certificate does not exists")
        }
    }

    if anode["type"] == "socket-network"{
        err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "interface", anode["interface"]); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "port", anode["port"]); if err != nil {logs.Error("InsertPluginService port Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "cert", anode["cert"]); if err != nil {logs.Error("InsertPluginService certtificate Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pid", "none"); if err != nil {logs.Error("InsertPluginService pid Error: "+err.Error()); return err}
    }
    if anode["type"] == "socket-pcap"{
        err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "interface", anode["interface"]); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "port", anode["port"]); if err != nil {logs.Error("InsertPluginService port Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "cert", anode["cert"]); if err != nil {logs.Error("InsertPluginService certtificate Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pcap-path", anode["pcap-path"]); if err != nil {logs.Error("InsertPluginService pcap-path Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pcap-prefix", anode["pcap-prefix"]); if err != nil {logs.Error("InsertPluginService pcap-prefix Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "bpf", anode["bpf"]); if err != nil {logs.Error("InsertPluginService bpf Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pid", "none"); if err != nil {logs.Error("InsertPluginService pid Error: "+err.Error()); return err}
    }
    if anode["type"] == "network-socket"{
        err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "interface", anode["interface"]); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "port", anode["port"]); if err != nil {logs.Error("InsertPluginService port Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "cert", anode["cert"]); if err != nil {logs.Error("InsertPluginService certtificate Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "collector", anode["collector"]); if err != nil {logs.Error("InsertPluginService collector Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "bpf", anode["bpf"]); if err != nil {logs.Error("InsertPluginService bpf Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pid", "none"); if err != nil {logs.Error("InsertPluginService pid Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "tcpdump", "none"); if err != nil {logs.Error("InsertPluginService pid Error: "+err.Error()); return err}
    }
    if anode["type"] == "zeek"{
        allPlugins,err := ndb.GetPlugins()
        for x := range allPlugins{
            if allPlugins[x]["type"] == "zeek"{ return errors.New("Can't Create more than one Zeek service.")}
        }
        err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "interface", "none"); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "status", "disabled"); if err != nil {logs.Error("InsertPluginService status Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "previousStatus", "none"); if err != nil {logs.Error("InsertPluginService previousStatus Error: "+err.Error()); return err}
    }
    if anode["type"] == "suricata"{
        path, err := utils.GetKeyValueString("suricataBPF", "pathBPF") 
        if err != nil {logs.Error("AddPluginService Error getting data from main.conf: "+err.Error())}

        // path := "/etc/suricata/bpf"
        if _, err := os.Stat(path); os.IsNotExist(err) { 
            err = os.MkdirAll(path, 0755); if err != nil {logs.Error("InsertPluginService erro creating BPF directory: "+err.Error()); return err}
        }

        err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "status", "disabled"); if err != nil {logs.Error("InsertPluginService status Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "previousStatus", "none"); if err != nil {logs.Error("InsertPluginService previousStatus Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "interface", ""); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "bpf", ""); if err != nil {logs.Error("InsertPluginService bpf Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "ruleset", ""); if err != nil {logs.Error("InsertPluginService ruleset Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pid", "none"); if err != nil {logs.Error("InsertPluginService ruleset Error: "+err.Error()); return err}
    }

    return nil
}

func SaveSuricataInterface(anode map[string]string)(err error) {
    err = ndb.UpdatePluginValue(anode["service"],anode["param"],anode["interface"])
    if err != nil {logs.Error("plugin/SaveSuricataInterface error: "+err.Error()); return err}
    return err
}

func CheckServicesStatus()(){
    plugin, err := utils.GetKeyValueString("stap", "plugin")
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    suriPID, err := utils.GetKeyValueString("execute", "suriPID")
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    param, err := utils.GetKeyValueString("execute", "param")
    if err != nil { logs.Error("Error getting data from main.conf")}
    command, err := utils.GetKeyValueString("execute", "command")
    if err != nil { logs.Error("Error getting data from main.conf")}
    socNetExec, err := utils.GetKeyValueString("execute", "socNetExec")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    socNetPID, err := utils.GetKeyValueString("execute", "socNetPID")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    socatPID, err := utils.GetKeyValueString("execute", "socatPID")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    openSSLparam, err := utils.GetKeyValueString("execute", "openSSLparam")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    tcpdump, err := utils.GetKeyValueString("execute", "tcpdump")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    socNetFile, err := utils.GetKeyValueString("execute", "socNetFile")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}

    allPlugins,_ := ndb.GetPlugins()
    for w := range allPlugins {
        if allPlugins[w]["pid"] != "none"{
            if allPlugins[w]["type"] == "suricata" {
                // pid, err := exec.Command("bash","-c","ps -ef | grep suricata | grep "+w+" | grep -v grep | awk '{print $2}'").Output()
                pid, err := exec.Command(command, param, strings.Replace(suriPID, "<ID>", w, -1)).Output()
                if err != nil {logs.Error("plugin/CheckServicesStatus Checking previous PID: "+err.Error())}
                pidValue := strings.Split(string(pid), "\n")
                
                if pidValue[0]!="" && pidValue[0] != allPlugins[w]["pid"] && allPlugins[w]["status"] == "enabled"{                    
                    err = ndb.UpdatePluginValue(w,"pid",pidValue[0])
                    if err != nil {logs.Error("plugin/CheckServicesStatus error updating pid at DB: "+err.Error())}
                    logs.Notice(pidValue[0]+" UPDATED!")
                }else if pidValue[0] == "" && allPlugins[w]["status"] == "enabled"{
                    err = LaunchSuricataService(w, allPlugins[w]["interface"])
                    if err != nil {
                        logs.Error("plugin/CheckServicesStatus error launching SURICATA after node stops: "+err.Error())
                        _ = StopSuricataService(w, allPlugins[w]["status"])
                    }else{
                        logs.Notice("Launching Suricata Service") 
                    } 
                }
            }else if allPlugins[w]["type"] == "zeek"{
                zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")  
                if err != nil {logs.Error("StopZeek Error getting data from main.conf: "+err.Error())}
                status, err := utils.GetKeyValueString("execute", "status")  
                if err != nil {logs.Error("StopZeek Error getting data from main.conf: "+err.Error())}

                // pid, err := exec.Command("bash","-c",zeekctl+" status | grep running | awk '{print $5}'").Output()
                pid, err := exec.Command("bash","-c",zeekctl+" "+status).Output()
                if err != nil {logs.Error("plugin/CheckServicesStatus Checking Zeek PID: "+err.Error())}
                
                if allPlugins[w]["status"] == "enabled"{                    
                    if (len(pid) == 0){
                        err = zeek.DeployZeek()                        
                        if err != nil {logs.Error("plugin/CheckQServicesStatus error deploying zeek: "+err.Error())}
                        // err = ndb.UpdatePluginValue(w,"pid",string(pid))
                        logs.Notice("Launch Zeek after Node stops")
                    }
                    // else{
                    //     if (allPlugins[w]["pid"] != string(pid)){
                    //         logs.Info("Zeek updated after Node stops")
                    //     }
                    // }
                }else if (allPlugins[w]["status"] == "disabled") {
                    if (len(pid) != 0){
                        _,err = zeek.StopZeek()
                        if err != nil {logs.Error("plugin/CheckServicesStatus error stopping zeek: "+err.Error())}
                        // err = ndb.UpdatePluginValue(w,"pid","none")                        
                        logs.Notice("Zeek stopped...")
                    }
                }
            }else if allPlugins[w]["type"] == "socket-network" {
                if allPlugins[w]["pid"] != "none" {   
                    // pid, err := exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[w]["port"]+" | grep -v grep | awk '{print $2}'").Output()
                    pid, err := exec.Command(command, param, strings.Replace(socatPID, "<PORT>", allPlugins[w]["port"], -1)).Output()
                    if err != nil {logs.Error("plugin/CheckServicesStatus Checking previous PID for socket-network: "+err.Error())}
                    pidValue := strings.Split(string(pid), "\n")
                    
                    if pidValue[0] == ""{
                        port := strings.Replace(socNetExec, "<PORT>", allPlugins[w]["port"], -1)
                        cert := strings.Replace(port, "<CERT>", allPlugins[w]["cert"], -1)
                        allValues := strings.Replace(cert, "<IFACE>", allPlugins[w]["interface"], -1)

                        // cmd := exec.Command("bash","-c",plugin+" -d OPENSSL-LISTEN:"+allPlugins[w]["port"]+",reuseaddr,pf=ip4,fork,cert="+allPlugins[w]["cert"]+",verify=0 SYSTEM:\"tcpreplay -t -i "+allPlugins[w]["interface"]+" -\" &")
                        cmd := exec.Command(command, param, plugin+" "+allValues)
                        var errores bytes.Buffer
                        cmd.Stdout = &errores
                        err = cmd.Start()
                        if err != nil {logs.Error("CheckServicesStatus deploying Error socket-network: "+err.Error())}        
                        
                        // pid, err = exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[w]["port"]+" | grep -v grep | awk '{print $2}'").Output()
                        pid, err = exec.Command(command, param, strings.Replace(socNetPID, "<PORT>", allPlugins[w]["port"], -1)).Output()
                        if err != nil {logs.Error("CheckServicesStatus deploy socket-network Error: "+err.Error())}
                        pidValue = strings.Split(string(pid), "\n")
                        if pidValue[0] != "" {
                            err = ndb.UpdatePluginValue(w,"pid",pidValue[0]); if err != nil {logs.Error("CheckServicesStatus change pid to value Error socket-network: "+err.Error())}
                        }
                    }
                }
            }else if  allPlugins[w]["type"] == "socket-pcap"{
                if allPlugins[w]["pid"] != "none" {                  
                    // pid, err := exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[w]["port"]+" | grep -v grep | awk '{print $2}'").Output()
                    pid, err := exec.Command(command, param, strings.Replace(socNetPID, "<PORT>", allPlugins[w]["port"], -1)).Output()
                    if err != nil {logs.Error("plugin/CheckServicesStatus Checking previous PID for socket-pcap: "+err.Error())}
                    pidValue := strings.Split(string(pid), "\n")
                    
                    if pidValue[0] == ""{
                        port := strings.Replace(socNetFile, "<PORT>", allPlugins[w]["port"], -1)
                        cert := strings.Replace(port, "<CERT>", allPlugins[w]["cert"], -1)
                        path := strings.Replace(cert, "<PCAP_PATH>", allPlugins[w]["pcap-path"], -1)
                        prefix := strings.Replace(path, "<PCAP_PREFIX>", allPlugins[w]["pcap-prefix"], -1)
                        allValues := strings.Replace(prefix, "<BPF>", allPlugins[w]["bpf"], -1)

                        // cmd := exec.Command("bash","-c",plugin+" -d OPENSSL-LISTEN:"+allPlugins[w]["port"]+",reuseaddr,pf=ip4,fork,cert="+allPlugins[w]["cert"]+",verify=0 SYSTEM:\"tcpdump -n -r - -s 0 -G 50 -W 100 -w "+allPlugins[w]["pcap-path"]+allPlugins[w]["pcap-prefix"]+"%d%m%Y%H%M%S.pcap "+allPlugins[w]["bpf"]+"\" &")
                        cmd := exec.Command(command, param, plugin+" "+allValues)
                        var errores bytes.Buffer
                        cmd.Stdout = &errores
                        err = cmd.Start()
                        if err != nil {logs.Error("CheckServicesStatus deploying Error socket-pcap: "+err.Error())}        
    
                        // pid, err = exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[w]["port"]+" | grep -v grep | awk '{print $2}'").Output()
                        pid, err = exec.Command(command, param, strings.Replace(socNetPID, "<PORT>", allPlugins[w]["port"], -1)).Output()
                        if err != nil {logs.Error("CheckServicesStatus deploy socket-pcap Error: "+err.Error())}
                        pidValue = strings.Split(string(pid), "\n")
                        if pidValue[0] != "" {
                            err = ndb.UpdatePluginValue(w,"pid",pidValue[0]); if err != nil {logs.Error("CheckServicesStatus change pid to value Error socket-pcap: "+err.Error())}
                        }
                    }
                }

            }else if allPlugins[w]["type"] == "network-socket"{
                if allPlugins[w]["pid"] != "none" && allPlugins[w]["tcpdump"] != "none" {             
                    collector := strings.Replace(openSSLparam, "<COLLECTOR>", allPlugins[w]["collector"], -1)
                    allValues := strings.Replace(collector, "<PORT>", allPlugins[w]["port"], -1)

                    // pidSocat, err := exec.Command("bash","-c","ps -ef | grep OPENSSL:"+allPlugins[w]["collector"]+":"+allPlugins[w]["port"]+" | awk '{print $2}'").Output()
                    pidSocat, err := exec.Command(command, param, allValues).Output()
                    if err != nil {logs.Error("CheckServicesStatus get socat pid at network-socket Error: "+err.Error())}
                    // pidTcpdump, err := exec.Command("bash","-c","ps -ef | grep tcpdump | grep -v grep | awk '{print $2}'").Output()
                    pidTcpdump, err := exec.Command(command, param, tcpdump).Output()
                    if err != nil {logs.Error("CheckServicesStatus get tcpdump at network-socket Error: "+err.Error())}

                    isRunningProcess := false
                    socatValue := strings.Split(string(pidSocat), "\n")
                    tcpdumpValue := strings.Split(string(pidTcpdump), "\n")
                    for _,value1 := range socatValue {
                        if allPlugins[w]["pid"] == value1 {
                            for _,value2 := range tcpdumpValue {
                                if allPlugins[w]["tcpdump"] == value2 {
                                    isRunningProcess = true                                    
                                }
                            }                            
                        }
                    }

                    if !isRunningProcess {
                        anode := make(map[string]string)
                        for x,y := range allPlugins{
                            for y,_ := range y{
                                if x == w {
                                    anode[y] = allPlugins[x][y]
                                }
                            }
                        }
                        anode["service"] = w

                        err = StopStapService(anode)
                        if err != nil {logs.Error("CheckServicesStatus error stopping node: "+err.Error())}                                    
                        err = DeployStapService(anode)
                        if err != nil {logs.Error("CheckServicesStatus error launching node: "+err.Error())}                                    
                    }                   
                }
            }
        }
    }
}

func LaunchSuricataService(uuid string, iface string)(err error){
    suricataBackup, err := utils.GetKeyValueString("suricata", "backup")
    if err != nil {logs.Error("LaunchSuricataService Error getting data from main.conf: "+err.Error())}
    filter, err := utils.GetKeyValueString("suricata", "filter")
    if err != nil {logs.Error("LaunchSuricataService Error getting data from main.conf: "+err.Error())}
    pidFile, err := utils.GetKeyValueString("suricata", "pidfile")
    if err != nil {logs.Error("LaunchSuricataService Error getting data from main.conf: "+err.Error())}
    fullpidfile, err := utils.GetKeyValueString("suricata", "fullpidfile")
    if err != nil {logs.Error("LaunchSuricataService Error getting data from main.conf: "+err.Error())}
    suricata, err := utils.GetKeyValueString("suricata", "suricata")
    if err != nil {logs.Error("LaunchSuricataService Error getting data from main.conf: "+err.Error())}
    param, err := utils.GetKeyValueString("execute", "param")
    if err != nil {logs.Error("LaunchSuricataService Error getting data from main.conf: "+err.Error())}
    suricata_config, err := utils.GetKeyValueString("files", "suricata_config")
    if err != nil {logs.Error("LaunchSuricataService Error getting data from main.conf: "+err.Error())}

    mainConfData, err := ndb.GetMainconfData()
    if (mainConfData["suricata"]["status"] == "disabled"){ return nil }

    _ = os.Remove(suricataBackup+uuid+"-"+pidFile)
    // cmd := exec.Command("suricata", "-D", "-c", "/etc/suricata/suricata.yaml", "-i", iface, "-F", "/etc/suricata/bpf/"+uuid+"-filter.bpf" ,"--pidfile", "/var/run/suricata/"+uuid+"-pidfile.pid")
    cmd := exec.Command(suricata, "-D", param, suricata_config, "-i", iface, "-F", strings.Replace(filter, "<ID>", uuid, -1) ,"--pidfile", strings.Replace(fullpidfile, "<ID>", uuid, -1))
    var stdBuffer bytes.Buffer
    cmd.Stderr = &stdBuffer

    err = cmd.Run()
    if err != nil {
        logs.Error(stdBuffer.String())
        logs.Error("plugin/LaunchSuricataService error launching Suricata: "+err.Error());
        //delete pid file
        err = os.Remove(suricataBackup+uuid+"-"+pidFile)
        if err != nil {logs.Error("plugin/LaunchSuricataService error deleting a pid file: "+err.Error())}
    }else{
        //read file
        currentpid, err := os.Open(suricataBackup+uuid+"-"+pidFile)
        if err != nil {logs.Error("plugin/LaunchSuricataService error openning Suricata: "+err.Error()); return err}
        defer currentpid.Close()
        pid, err := ioutil.ReadAll(currentpid)
        dbValue := strings.Split(string(pid), "\n")

        //save pid to db
        err = ndb.UpdatePluginValue(uuid,"pid",dbValue[0])
        if err != nil {logs.Error("plugin/LaunchSuricataService error updating pid at DB: "+err.Error()); return err}

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
    suricataBackup, err := utils.GetKeyValueString("suricata", "backup")  
    if err != nil {logs.Error("StopSuricataService Error getting data from main.conf: "+err.Error())}
    suricataPidfile, err := utils.GetKeyValueString("suricata", "pidfile")  
    if err != nil {logs.Error("StopSuricataService Error getting data from main.conf: "+err.Error())}
    //pid
    // currentpid, err := os.Open("/var/run/suricata/"+uuid+"-pidfile.pid")
    // if err != nil {logs.Error("plugin/ChangeServiceStatus error reading Suricata pid: "+err.Error()); return err}
    // defer currentpid.Close()
    // pid, err := ioutil.ReadAll(currentpid)
    allPlugins,err := ndb.GetPlugins()

    //kill suricata process
    PidInt,_ := strconv.Atoi(strings.Trim(string(allPlugins[uuid]["pid"]), "\n"))
    process, _ := os.FindProcess(PidInt)
    _ = process.Kill()
    // if err != nil {logs.Error("plugin/StopSuricataService error killing Suricata process: "+err.Error()); return err}

    //delete pid file
    _ = os.Remove(suricataBackup+uuid+"-"+suricataPidfile)
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

func ModifyStapValues(anode map[string]string)(err error) {
    allPlugins,err := ndb.GetPlugins()
    if anode["type"] == "zeek"{
        err = ndb.UpdatePluginValue(anode["service"],"name",anode["name"]); if err != nil {logs.Error("ModifyStapValues zeek Error: "+err.Error()); return err}
        if allPlugins[anode["service"]]["status"] == "enabled" {
            err = zeek.DeployZeek()
            if err != nil {logs.Error("plugin/ModifyStapValues error deploying zeek: "+err.Error()); return err}
        }
        logs.Notice(allPlugins[anode["service"]]["name"]+" service updated!!!")
    }else if anode["type"] == "suricata"{
        err = ndb.UpdatePluginValue(anode["service"],"name",anode["name"]); if err != nil {logs.Error("ModifyStapValues suricata Error: "+err.Error()); return err}
        if allPlugins[anode["service"]]["status"] == "enabled" {
            err = StopSuricataService(anode["service"], allPlugins[anode["service"]]["status"])
            if err != nil {logs.Error("plugin/ModifyStapValues error stopping suricata: "+err.Error()); return err}
            err = LaunchSuricataService(anode["service"], allPlugins[anode["service"]]["interface"])
            if err != nil {logs.Error("plugin/ModifyStapValues error deploying suricata: "+err.Error()); return err}
        }
        logs.Notice(allPlugins[anode["service"]]["name"]+" service updated!!!")
    }else if anode["type"] == "socket-pcap" || anode["type"] == "socket-network"{
        err = ndb.UpdatePluginValue(anode["service"],"name",anode["name"]) ; if err != nil {logs.Error("ModifyStapValues "+anode["type"]+" Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"port",anode["port"]) ; if err != nil {logs.Error("ModifyStapValues "+anode["type"]+" Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"cert",anode["cert"]) ; if err != nil {logs.Error("ModifyStapValues "+anode["type"]+" Error: "+err.Error()); return err}
        if anode["type"] == "socket-pcap"{
            err = ndb.UpdatePluginValue(anode["service"],"pcap-path",anode["pcap-path"]) ; if err != nil {logs.Error("ModifyStapValues "+anode["type"]+" Error: "+err.Error()); return err}
            err = ndb.UpdatePluginValue(anode["service"],"pcap-prefix",anode["pcap-prefix"]) ; if err != nil {logs.Error("ModifyStapValues "+anode["type"]+" Error: "+err.Error()); return err}
        }
        for x := range allPlugins{
            if ((allPlugins[x]["type"] == "socket-network" || allPlugins[x]["type"] == "socket-pcap") && (anode["service"] != x)){
                if allPlugins[x]["port"] == anode["port"] {
                    err = StopStapService(anode); if err != nil {logs.Error("ModifyStapValues socket-network stopping error: "+err.Error()); return err}        
                    logs.Error("Can't deploy socket-network or "+anode["type"]+" with the same port")
                    return errors.New("Can't deploy socket-network or "+anode["type"]+" with the same port")
                }
            }
        }
        if allPlugins[anode["service"]]["pid"] != "none" {
            err = StopStapService(anode); if err != nil {logs.Error("ModifyStapValues "+anode["type"]+" stopping error: "+err.Error()); return err}
            err = DeployStapService(anode); if err != nil {logs.Error("ModifyStapValues "+anode["type"]+" deploying error: "+err.Error()); return err}
            logs.Notice(allPlugins[anode["service"]]["name"]+" service updated!!!")
        }
    }else if anode["type"] == "network-socket"{
        err = ndb.UpdatePluginValue(anode["service"],"name",anode["name"]) ; if err != nil {logs.Error("ModifyStapValues network-socket Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"port",anode["port"]) ; if err != nil {logs.Error("ModifyStapValues network-socket Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"cert",anode["cert"])  ; if err != nil {logs.Error("ModifyStapValues network-socket Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"collector",anode["collector"]) ; if err != nil {logs.Error("ModifyStapValues network-socket Error: "+err.Error()); return err}
        for x := range allPlugins{
            if x != anode["service"] && allPlugins[x]["type"] == anode["type"] && allPlugins[x]["collector"] == anode["collector"] && allPlugins[x]["port"] == anode["port"] && allPlugins[x]["interface"] == anode["interface"]{
                logs.Error("This network-socket has been deployed yet. Can't update")
                err = StopStapService(anode); if err != nil {logs.Error("ModifyStapValues error stopping duplicated network-socket: "+err.Error()); return err}
                return errors.New("This network-socket has been deployed yet. Can't update")
            }
        }        
        if allPlugins[anode["service"]]["pid"] != "none" && allPlugins[anode["service"]]["tcpdump"] != "none"{
            logs.Notice("Updating "+allPlugins[anode["service"]]["name"]+" service...")
            err = StopStapService(anode); if err != nil {logs.Error("ModifyStapValues network-socket stopping error: "+err.Error()); return err}
            err = DeployStapService(anode); if err != nil {logs.Error("ModifyStapValues network-socket deploying error: "+err.Error()); return err}
            logs.Notice(allPlugins[anode["service"]]["name"]+" service updated!!!")
        }
    }
    return nil
}

func DeployStapService(anode map[string]string)(err error) {     
    stapPlugin, err := utils.GetKeyValueString("stap", "plugin")
    if err != nil {logs.Error("DeployStapService Error getting data from main.conf: "+err.Error())}
    stapTcpdump, err := utils.GetKeyValueString("stap", "tcpdum")
    if err != nil {logs.Error("DeployStapService Error getting data from main.conf: "+err.Error())}
    param, err := utils.GetKeyValueString("execute", "param")
    if err != nil { logs.Error(" Error getting data from main.conf")}
    command, err := utils.GetKeyValueString("execute", "command")
    if err != nil { logs.Error("Error getting data from main.conf")}
    socNetPID, err := utils.GetKeyValueString("execute", "socNetPID")
    if err != nil { logs.Error("Error getting data from main.conf")}
    socNetExec, err := utils.GetKeyValueString("execute", "socNetExec")
    if err != nil { logs.Error("Error getting data from main.conf")}
    socNetFile, err := utils.GetKeyValueString("execute", "socNetFile")
    if err != nil { logs.Error("Error getting data from main.conf")}
    NetSocFile, err := utils.GetKeyValueString("execute", "NetSocFile")
    if err != nil { logs.Error("Error getting data from main.conf")}
    openSSLparam, err := utils.GetKeyValueString("execute", "openSSLparam")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    tcpdumpPID, err := utils.GetKeyValueString("execute", "tcpdumpPID")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}

    allPlugins,err := ndb.GetPlugins()
    if anode["type"] == "socket-network" {
        // pid, err := exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+" | grep -v grep | awk '{print $2}'").Output()
        pid, err := exec.Command(command, param, strings.Replace(socNetPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)).Output()
        if err != nil {logs.Error("DeployStapService deploy socket-network Error: "+err.Error()); return err}
        pidValue := strings.Split(string(pid), "\n")
        if pidValue[0] != "" {
            logs.Error("Socket to network deployed. Exiting DeployStapService")
            return errors.New("Can't deploy more than one socket at the same port")            
        }

        port := strings.Replace(socNetExec, "<PORT>", allPlugins[anode["service"]]["port"], -1)
        cert := strings.Replace(port, "<CERT>", allPlugins[anode["service"]]["cert"], -1)
        allValues := strings.Replace(cert, "<IFACE>", allPlugins[anode["service"]]["interface"], -1)
        // cmd := exec.Command("bash","-c","/usr/bin/socat -d OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+",reuseaddr,pf=ip4,fork,cert="+allPlugins[anode["service"]]["cert"]+",verify=0 SYSTEM:\"tcpreplay -t -i "+allPlugins[anode["service"]]["interface"]+" -\" &")
        cmd := exec.Command(command, param, stapPlugin+" "+allValues)
        var errores bytes.Buffer
        cmd.Stdout = &errores
        err = cmd.Start()
        if err != nil {logs.Error("DeployStapService deploying Error: "+err.Error()); return err}        

        // pid, err = exec.Command(command, param, "ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+" | grep -v grep | awk '{print $2}'").Output()
        pid, err = exec.Command(command, param, strings.Replace(socNetPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)).Output()
        if err != nil {logs.Error("DeployStapService deploy socket-network Error: "+err.Error()); return err}
        pidValue = strings.Split(string(pid), "\n")
        if pidValue[0] != "" {
            err = ndb.UpdatePluginValue(anode["service"],"pid",pidValue[0]); if err != nil {logs.Error("DeployStapService change pid to value Error: "+err.Error()); return err}
        }
        logs.Notice("Deploy successful --> Type: "+allPlugins[anode["service"]]["type"]+" Description: "+allPlugins[anode["service"]]["name"]+"  --  SOCAT: "+pidValue[0])
    }else if anode["type"] == "socket-pcap" {
        // pid, err := exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+" | grep -v grep | awk '{print $2}'").Output()
        pid, err := exec.Command(command, param, strings.Replace(socNetPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)).Output()
        if err != nil {logs.Error("DeployStapService deploy socket-network Error: "+err.Error()); return err}
        pidValue := strings.Split(string(pid), "\n")
        if pidValue[0] != "" {
            logs.Error("Socket to pcap deployed. Exiting DeployStapService")
            return errors.New("Can't deploy more than one socket at the same port")            
        }

        // cmd := exec.Command("bash","-c",stapPlugin+" -d OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+",reuseaddr,pf=ip4,fork,cert="+allPlugins[anode["service"]]["cert"]+",verify=0 SYSTEM:\"tcpdump -n -r - -s 0 -G 50 -W 100 -w "+allPlugins[anode["service"]]["pcap-path"]+allPlugins[anode["service"]]["pcap-prefix"]+"%d%m%Y%H%M%S.pcap "+allPlugins[anode["service"]]["bpf"]+"\" &")
        port := strings.Replace(socNetFile, "<PORT>",allPlugins[anode["service"]]["port"], -1)
        cert := strings.Replace(port, "<CERT>", allPlugins[anode["service"]]["cert"], -1)
        path := strings.Replace(cert, "<PCAP_PATH>",allPlugins[anode["service"]]["pcap-path"], -1)
        prefix := strings.Replace(path, "<PCAP_PREFIX>", allPlugins[anode["service"]]["pcap-prefix"], -1)
        allValues := strings.Replace(prefix, "<BPF>", allPlugins[anode["service"]]["bpf"], -1)

        // cmd := exec.Command("bash","-c",stapPlugin+" -d OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+",reuseaddr,pf=ip4,fork,cert="+allPlugins[anode["service"]]["cert"]+",verify=0 SYSTEM:\"tcpdump -n -r - -s 0 -G 50 -W 100 -w "+allPlugins[anode["service"]]["pcap-path"]+allPlugins[anode["service"]]["pcap-prefix"]+"%d%m%Y%H%M%S.pcap "+allPlugins[anode["service"]]["bpf"]+"\" &")
        cmd := exec.Command(command, param, stapPlugin+" "+allValues )
        var errores bytes.Buffer
        cmd.Stdout = &errores
        err = cmd.Start()
        if err != nil {logs.Error("DeployStapService deploying Error: "+err.Error()); return err}        

        // pid, err = exec.Command(command, param, "ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+" | grep -v grep | awk '{print $2}'").Output()
        pid, err = exec.Command(command, param, strings.Replace(socNetPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)).Output()
        if err != nil {logs.Error("DeployStapService deploy socket-network Error: "+err.Error()); return err}
        pidValue = strings.Split(string(pid), "\n")
        if pidValue[0] != "" {
            err = ndb.UpdatePluginValue(anode["service"],"pid",pidValue[0]); if err != nil {logs.Error("DeployStapService change pid to value Error: "+err.Error()); return err}
        }
        logs.Notice("Deploy successful --> Type: "+allPlugins[anode["service"]]["type"]+" Description: "+allPlugins[anode["service"]]["name"]+"  --  SOCAT: "+pidValue[0])
    }else if anode["type"] == "network-socket" {
        for x := range allPlugins{
            if x != anode["service"] && allPlugins[x]["type"] == anode["type"] && allPlugins[x]["collector"] == anode["collector"] && allPlugins[x]["port"] == anode["port"] && allPlugins[x]["interface"] == anode["interface"]{
                logs.Error("This network-socket has been deployed yet")
                return errors.New("This network-socket has been deployed yet")
            }
        }

        iface := strings.Replace(NetSocFile, "<IFACE>", allPlugins[anode["service"]]["interface"], -1)
        bpf := strings.Replace(iface, "<BPF>", allPlugins[anode["service"]]["bpf"], -1)
        stap := strings.Replace(bpf, "<STAP>", stapPlugin, -1)
        collector := strings.Replace(stap, "<COLLECTOR>", allPlugins[anode["service"]]["collector"], -1)
        port := strings.Replace(collector, "<PORT>", allPlugins[anode["service"]]["port"], -1)
        allNetSock := strings.Replace(port, "<CERT>", allPlugins[anode["service"]]["cert"], -1)

        // cmd := exec.Command("bash","-c",stapTcpdump+" -n -i "+allPlugins[anode["service"]]["interface"]+" -s 0 -w - "+allPlugins[anode["service"]]["bpf"]+" | "+stapPlugin+" - OPENSSL:"+allPlugins[anode["service"]]["collector"]+":"+allPlugins[anode["service"]]["port"]+",cert="+allPlugins[anode["service"]]["cert"]+",verify=0,forever,retry=10,interval=5 &")
        cmd := exec.Command(command, param, stapTcpdump+" "+allNetSock)
        err = cmd.Start()
        if err != nil {logs.Error("DeployStapService deploying Error: "+err.Error()); return err}

        //get socat pid
        var grepPIDS string
        for x := range allPlugins{
            if allPlugins[x]["type"] == "network-socket" && allPlugins[x]["pid"] != "none"{
                grepPIDS = grepPIDS + "| grep -v "+allPlugins[x]["pid"]+" "
            }
        }

        collector = strings.Replace(openSSLparam, "<COLLECTOR>", allPlugins[anode["service"]]["collector"], -1)
        allValues := strings.Replace(collector, "<PORT>", allPlugins[anode["service"]]["port"], -1)
        // pid, err := exec.Command("bash","-c","ps -ef | grep OPENSSL:"+allPlugins[anode["service"]]["collector"]+":"+allPlugins[anode["service"]]["port"]+" "+grepPIDS+" | grep -v grep | awk '{print $2}'").Output()
        pid, err := exec.Command(command, param, allValues).Output()
        if err != nil {logs.Error("DeployStapService deploy network-socket getting socat error: "+err.Error()); return err}
        pidValueSocat := strings.Split(string(pid), "\n")

        if pidValueSocat[0] != "" {
            err = ndb.UpdatePluginValue(anode["service"],"pid",pidValueSocat[0]); if err != nil {logs.Error("DeployStapService update DB pid to value Error: "+err.Error()); return err}
        }

        //get tcpdump
        var grepTCPDUMP string
        for x := range allPlugins{
            if allPlugins[x]["type"] == "network-socket" && allPlugins[x]["tcpdump"] != "none"{
                grepTCPDUMP = grepTCPDUMP + "| grep -v "+allPlugins[x]["tcpdump"]+" "
            }
        }
        tcpdump := strings.Replace(tcpdumpPID, "<TCPDUMP>", grepTCPDUMP, -1)
        // pid, err = exec.Command("bash","-c","ps -ef | grep -v grep | grep tcpdump "+grepTCPDUMP+" | awk '{print $2}'").Output()
        pid, err = exec.Command(command, param, tcpdump).Output()
        if err != nil {logs.Error("DeployStapService deploy network-socket getting tcpdump pid error: "+err.Error()); return err}
        pidValueTcpdump := strings.Split(string(pid), "\n")

        if pidValueTcpdump[0] != "" {
            err = ndb.UpdatePluginValue(anode["service"],"tcpdump",pidValueTcpdump[0]); if err != nil {logs.Error("DeployStapService update DB tcpdump to value Error: "+err.Error()); return err}
        }

        logs.Notice("Deploy successful --> Type: "+allPlugins[anode["service"]]["type"]+" Description: "+allPlugins[anode["service"]]["name"]+"  --  SOCAT: "+pidValueSocat[0]+"  --  TCPDUMP: "+pidValueTcpdump[0])
    }
    
    return nil
}

func StopStapService(anode map[string]string)(err error) {
    logs.Debug(anode)
    allPlugins,err := ndb.GetPlugins()
    if err != nil {logs.Error("Error! can't read database for stop the service: "+err.Error())}
    pidToInt,_ := strconv.Atoi(allPlugins[anode["service"]]["pid"])
    process, _ := os.FindProcess(pidToInt)
    _ = process.Kill()
    err = ndb.UpdatePluginValue(anode["service"],"pid","none") ; if err != nil {logs.Error("DeployStapService update DB pid to none Error: "+err.Error()); return err}

    if allPlugins[anode["service"]]["type"] == "network-socket" {
        //kill tcpdump
        tcpdumpToInt,_ := strconv.Atoi(allPlugins[anode["service"]]["tcpdump"])
        processTcpdump, _ := os.FindProcess(tcpdumpToInt)
        _ = processTcpdump.Kill()
        err = ndb.UpdatePluginValue(anode["service"],"tcpdump","none") ; if err != nil {logs.Error("DeployStapService update DB tcpdump to none Error: "+err.Error()); return err}
    }
    logs.Notice(allPlugins[anode["service"]]["type"]+" service stopped successfuly!")

    return nil
}

func ChangeSuricataTable(anode map[string]string)(err error) {
    // allPlugins,err := ndb.GetPlugins()
    data, err := ndb.GetMainconfData()

    if anode["status"] == "expert" {
        err = ndb.UpdateMainconfValue("suricata", "previousStatus", data["suricata"]["status"]); if err != nil {logs.Error("ChangeSuricataTable status Error: "+err.Error()); return err}
        err = ndb.UpdateMainconfValue("suricata", "status", "expert"); if err != nil {logs.Error("ChangeSuricataTable status Error: "+err.Error()); return err}
    }else{
        if data["suricata"]["previousStatus"] == "enabled" {
            err = ndb.UpdateMainconfValue("suricata", "status", data["suricata"]["previousStatus"])
            err = ndb.UpdateMainconfValue("suricata", "previousStatus", "disabled")  
        }else if data["suricata"]["previousStatus"] == "disabled"{
            err = ndb.UpdateMainconfValue("suricata", "status", data["suricata"]["previousStatus"])
            err = ndb.UpdateMainconfValue("suricata", "previousStatus", "enabled")
        }else {
            ndb.InsertGetMainconfData("suricata", "previousStatus", "disabled")
        }
    } 

    // for x := range allPlugins {
    //     if anode["status"] == "expert" {
    //         if allPlugins[x]["status"] == "enabled" && allPlugins[x]["type"] == "suricata"{
    //             err = StopSuricataService(x, allPlugins[x]["status"])
    //             if err != nil {logs.Error("StopSuricataService status Error: "+err.Error()); return err}
    //         } 
    //     }else{
    //         if data["suricata"]["previousStatus"] == "enabled" {
    //             if allPlugins[x]["previousState"] == "enabled" && allPlugins[x]["type"] == "suricata"{
    //                 err = LaunchSuricataService(x, allPlugins[x]["interface"])
    //                 if err != nil {logs.Error("LaunchSuricataService status Error: "+err.Error()); return err}
    //             }          
    //         }else if data["suricata"]["previousStatus"] == "disabled"{
    //             if allPlugins[x]["previousStatus"] == "enabled" && allPlugins[x]["type"] == "suricata"{
    //                 err = StopSuricataService(x, allPlugins[x]["status"])
    //                 if err != nil {logs.Error("ChangeSuricataTable LaunchSuricataService status Error: "+err.Error()); return err}
    //             }
    //         }else {
    //             ndb.InsertGetMainconfData("suricata", "previousStatus", "disabled")
    //         }
    //     } 
    // }
    
    return nil
}