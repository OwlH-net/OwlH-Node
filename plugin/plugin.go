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
    "time"
    "strconv"
    "strings"
    "io/ioutil"
    "owlhnode/utils"
)

func ChangeServiceStatus(anode map[string]string)(err error) {
    allPlugins,err := ndb.GetPlugins()
    if anode["type"] == "suricata"{
        mainConfData, err := ndb.GetMainconfData()
        if (mainConfData["suricata"]["status"] == "disabled"){ return errors.New("Please, enable main Suricata status before launching") }    
        if anode["status"] == "enabled"{
            for x := range allPlugins {
                //check if suricata is deployed yet
                if allPlugins[x]["type"] == "suricata" && x == anode["service"] && allPlugins[x]["status"] == "enabled"{
                    return nil
                }
                //get all db values and check if there are any suricata deployed at the same interface
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
        if (mainConfData["zeek"]["status"] == "disabled"){ return errors.New("Please, enable main Zeek status before launch") }        
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
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "interface", anode["interface"]); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "port", anode["port"]); if err != nil {logs.Error("InsertPluginService port Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "cert", anode["cert"]); if err != nil {logs.Error("InsertPluginService certtificate Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pid", "none"); if err != nil {logs.Error("InsertPluginService pid Error: "+err.Error()); return err}
    }
    if anode["type"] == "socket-pcap"{
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
    suriPID, err := utils.GetKeyValueString("execute", "suriPID")
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    param, err := utils.GetKeyValueString("execute", "param")
    if err != nil { logs.Error("Error getting data from main.conf")}
    command, err := utils.GetKeyValueString("execute", "command")
    if err != nil { logs.Error("Error getting data from main.conf")}

    time.Sleep(time.Second * 2) 
    allPlugins,_ := ndb.GetPlugins()
    for w := range allPlugins {
        if allPlugins[w]["pid"] != "none"{
            if allPlugins[w]["type"] == "suricata" {
                pid, err := exec.Command(command, param, strings.Replace(suriPID, "<ID>", "grep "+w+" |", -1)).Output()
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

                pid, err := exec.Command("bash","-c",zeekctl+" "+status).Output()
                if err != nil {logs.Error("plugin/CheckServicesStatus Checking Zeek PID: "+err.Error())}

                if allPlugins[w]["status"] == "enabled"{                    
                    if (len(pid) == 0){
                        err = zeek.DeployZeek()                        
                        if err != nil {logs.Error("plugin/CheckQServicesStatus error deploying zeek: "+err.Error())}
                        logs.Notice("Launch Zeek after Node stops")
                    }
                }else if (allPlugins[w]["status"] == "disabled") {
                    if (len(pid) != 0){
                        _,err = zeek.StopZeek()
                        if err != nil {logs.Error("plugin/CheckServicesStatus error stopping zeek: "+err.Error())}
                        logs.Notice("Zeek stopped...")
                    }
                }
            }else if allPlugins[w]["type"] == "socket-network" {
                if allPlugins[w]["pid"] != "none" {   
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
            }else if  allPlugins[w]["type"] == "socket-pcap"{
                if allPlugins[w]["pid"] != "none" {   
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

            }else if allPlugins[w]["type"] == "network-socket"{
                if allPlugins[w]["pid"] != "none" && allPlugins[w]["tcpdump"] != "none" {             
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
        time.Sleep(time.Second * 1)
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
    allPlugins,err := ndb.GetPlugins()

    //kill suricata process
    PidInt,_ := strconv.Atoi(strings.Trim(string(allPlugins[uuid]["pid"]), "\n"))
    process, _ := os.FindProcess(PidInt)
    _ = process.Kill()

    //delete pid file
    _ = os.Remove(suricataBackup+uuid+"-"+suricataPidfile)

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

        //check for STAP certificate
        if _, err := os.Stat(anode["cert"]); os.IsNotExist(err) {
            logs.Error("STAP certificate does not exists")
            err = StopStapService(anode); if err != nil {logs.Error("ModifyStapValues socket-network stopping error: "+err.Error()); return err}
            return errors.New("STAP certificate does not exists")
        }   

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
        //check for STAP certificate
        if _, err := os.Stat(anode["cert"]); os.IsNotExist(err) {
            logs.Error("STAP certificate does not exists")
            err = StopStapService(anode); if err != nil {logs.Error("ModifyStapValues socket-network stopping error: "+err.Error()); return err}
            return errors.New("STAP certificate does not exists")
        }   

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
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    stapTcpdump, err := utils.GetKeyValueString("stap", "tcpdum")
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    command, err := utils.GetKeyValueString("execute", "command")
    if err != nil { logs.Error("Error getting data from main.conf")}
    param, err := utils.GetKeyValueString("execute", "param")
    if err != nil { logs.Error(" Error getting data from main.conf")}
    socatPID, err := utils.GetKeyValueString("execute", "socatPID")
    if err != nil { logs.Error("Error getting data from main.conf")}
    socNetExec, err := utils.GetKeyValueString("execute", "socNetExec")
    if err != nil { logs.Error("Error getting data from main.conf")}
    socNetFile, err := utils.GetKeyValueString("execute", "socNetFile")
    if err != nil { logs.Error("Error getting data from main.conf")}
    NetSocFile, err := utils.GetKeyValueString("execute", "NetSocFile")
    if err != nil { logs.Error("Error getting data from main.conf")}
    openSSL, err := utils.GetKeyValueString("execute", "openSSL")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    tcpdumpPID, err := utils.GetKeyValueString("execute", "tcpdumpPID")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}

    //insert common values into command db
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "id", anode["service"])
    _ = ndb.InsertPluginCommand(uuid, "type", anode["type"])
    _ = ndb.InsertPluginCommand(uuid, "action", "Deploy")

    allPlugins,err := ndb.GetPlugins()
    if anode["type"] == "socket-network" {
        //values for deploy socat-network
        port := strings.Replace(socNetExec, "<PORT>", allPlugins[anode["service"]]["port"], -1)
        cert := strings.Replace(port, "<CERT>", allPlugins[anode["service"]]["cert"], -1)
        allValues := strings.Replace(cert, "<IFACE>", allPlugins[anode["service"]]["interface"], -1)

        //insert command values
        _ = ndb.InsertPluginCommand(uuid, "command", stapPlugin+" "+allValues)
        _ = ndb.InsertPluginCommand(uuid, "description", "Deploy socket-network")

        //check if a selected STAP server is deployed yet
        pid, err := exec.Command(command, param, strings.Replace(socatPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)).Output()
        if err != nil {
            logs.Error("DeployStapService get socket-network PID Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", "Cannot check if a socket->network service has been launched yet at port "+allPlugins[anode["service"]]["port"])
            return errors.New("Cannot check if a socket->network service has been launched yet at port "+allPlugins[anode["service"]]["port"]+".")
        }
        pidValue := strings.Split(string(pid), "\n")       
        if pidValue[0] != "" {
            logs.Error("Socket to network running at port "+allPlugins[anode["service"]]["port"]+". Exiting DeployStapService")
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", "There is already a socket->network running at port "+allPlugins[anode["service"]]["port"])
            return errors.New("Cannot deploy more than one socket at the same port")
        }
        
        //deploy socat 
        cmd := exec.Command(command, param, stapPlugin+" "+allValues)
        stdError,err := cmd.StderrPipe()
        err = cmd.Start()
        if err != nil {logs.Error("CMD START ERROR --> "+ err.Error())}
        // logs.Warn(cmd.Process.Pid)

        time.Sleep(time.Second * 2) 
        //Get deployed PID
        pid, err = exec.Command(command, param, strings.Replace(socatPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)).Output()
        if err != nil {
            logs.Error("DeployStapService deploy socket-network Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", "Error deploying Socat service.")
            return err
        }
        pidValue = strings.Split(string(pid), "\n")
        
        if pidValue[0] == ""{
            pipeError, _ := ioutil.ReadAll(stdError)
            logs.Error(string(pipeError))
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", string(pipeError))
            return errors.New("Deploy socket-network error: Please, check warning log")
        }else{
            //update service status
            err = ndb.UpdatePluginValue(anode["service"],"pid",pidValue[0])
            if err != nil {
                logs.Error("DeployStapService change pid to value Error: "+err.Error())
                _ = ndb.InsertPluginCommand(uuid, "status", "Error")
                _ = ndb.InsertPluginCommand(uuid, "output", "Error updating Socat PID after deploy.")
                return err
            }
            _ = ndb.InsertPluginCommand(uuid, "output", pidValue[0])
            _ = ndb.InsertPluginCommand(uuid, "status", "Success")
        }
        logs.Notice("Deploy "+allPlugins[anode["service"]]["type"]+" successfully -->  Description: "+allPlugins[anode["service"]]["name"]+"  --  SOCAT: "+pidValue[0])
    }else if anode["type"] == "socket-pcap" {       
        //get socat command
        port := strings.Replace(socNetFile, "<PORT>",allPlugins[anode["service"]]["port"], -1)
        cert := strings.Replace(port, "<CERT>", allPlugins[anode["service"]]["cert"], -1)
        path := strings.Replace(cert, "<PCAP_PATH>",allPlugins[anode["service"]]["pcap-path"], -1)
        prefix := strings.Replace(path, "<PCAP_PREFIX>", allPlugins[anode["service"]]["pcap-prefix"], -1)
        allValues := strings.Replace(prefix, "<BPF>", allPlugins[anode["service"]]["bpf"], -1)
        
        //add Command Values into database
        _ = ndb.InsertPluginCommand(uuid, "command", stapPlugin+" "+allValues)
        _ = ndb.InsertPluginCommand(uuid, "description", "Deploy socket-pcap")
        
        //Check if a socket-pcap is deployed
        pid, err := exec.Command(command, param, strings.Replace(socatPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)).Output()
        if err != nil {
            logs.Error("DeployStapService deploy socket-pcap Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", "Error checking if Socat service is already deployed.")
            return err
        }
        pidValue := strings.Split(string(pid), "\n")
        if pidValue[0] != "" {
            logs.Error("Socket to pcap deployed. Exiting DeployStapService")
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", "There is already a socket->pcap running at port "+allPlugins[anode["service"]]["port"]+".")
            return errors.New("Can't deploy more than one socket at the same port")            
        }

        //deploy socat for socket->pcap
        cmd := exec.Command(command, param, stapPlugin+" "+allValues)
        stdError,err := cmd.StderrPipe()
        err = cmd.Start()
        if err != nil {
            logs.Error("DeployStapService deploying Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", "Error deploying Socat service for socket->pcap.")
            return err
        }

        time.Sleep(time.Second * 2) 
        //get PID
        pid, err = exec.Command(command, param, strings.Replace(socatPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)).Output()
        if err != nil {
            logs.Error("DeployStapService getting socket-network PID Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", "Error getting new Socat PID for socket->pcap service.")
            return err
        }
        pidValue = strings.Split(string(pid), "\n")

        //Check for the new socat PID
        if pidValue[0] == ""{
            pipeError, _ := ioutil.ReadAll(stdError)
            logs.Error(string(pipeError))
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", string(pipeError))
            return errors.New("Deploy socket-pcap error: Please, check warning log")
        }else{
            //update service status if there are a PID
            err = ndb.UpdatePluginValue(anode["service"],"pid",pidValue[0])
            if err != nil {
                logs.Error("DeployStapService change pid to value Error: "+err.Error())
                _ = ndb.InsertPluginCommand(uuid, "status", "Error")
                _ = ndb.InsertPluginCommand(uuid, "output", "Error updating database PID value for socket->pcap service.")
                return err
            }
            _ = ndb.InsertPluginCommand(uuid, "output", pidValue[0])
            _ = ndb.InsertPluginCommand(uuid, "status", "Success")
        }

        logs.Notice("Deploy "+allPlugins[anode["service"]]["type"]+" successfully --> Description: "+allPlugins[anode["service"]]["name"]+"  --  SOCAT: "+pidValue[0])
    }else if anode["type"] == "network-socket" {
        //tcpdump values for network->socket
        iface := strings.Replace(NetSocFile, "<IFACE>", allPlugins[anode["service"]]["interface"], -1)
        bpf := strings.Replace(iface, "<BPF>", allPlugins[anode["service"]]["bpf"], -1)
        stap := strings.Replace(bpf, "<STAP>", stapPlugin, -1)
        collector := strings.Replace(stap, "<COLLECTOR>", allPlugins[anode["service"]]["collector"], -1)
        port := strings.Replace(collector, "<PORT>", allPlugins[anode["service"]]["port"], -1)
        allNetSock := strings.Replace(port, "<CERT>", allPlugins[anode["service"]]["cert"], -1)

        //add Command Values into database
        _ = ndb.InsertPluginCommand(uuid, "command", stapTcpdump+" "+allNetSock)
        _ = ndb.InsertPluginCommand(uuid, "description", "Deploy tcpdump for network->socket")

        //check for a network->socket deployed
        for x := range allPlugins{
            if x != anode["service"] && allPlugins[x]["type"] == anode["type"] && allPlugins[x]["collector"] == anode["collector"] &&
            allPlugins[x]["port"] == anode["port"] && allPlugins[x]["interface"] == anode["interface"] && allPlugins[x]["pid"] != "none" && allPlugins[x]["tcpdump"] != "none"{
                logs.Error("This network-socket has been deployed yet")
                _ = ndb.InsertPluginCommand(uuid, "status", "Error")
                _ = ndb.InsertPluginCommand(uuid, "output", "There is already a network->socket running at this port, collector and interface")
                return errors.New("This network-socket has been deployed yet")
            }
        }

        //exec tcpdump
        cmd := exec.Command(command, param, stapTcpdump+" "+allNetSock)
        stdError,err := cmd.StderrPipe()
        err = cmd.Start()
        if err != nil {
            logs.Error("DeployStapService deploying Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", "Error starting tcpdump for network->socket service")
            return err
        }

        //get tcpreplay pid
        var grepPIDS string
        for x := range allPlugins{
            if allPlugins[x]["type"] == "network-socket" && allPlugins[x]["pid"] != "none"{
                grepPIDS = grepPIDS + "| grep -v "+allPlugins[x]["pid"]+" "
            }
        }

        time.Sleep(time.Second * 1) 
        //get tcpreplay pids
        collector = strings.Replace(openSSL, "<COLLECTOR>", allPlugins[anode["service"]]["collector"], -1)
        allValues := strings.Replace(collector, "<PORT>", allPlugins[anode["service"]]["port"], -1)
        pid, err := exec.Command(command, param, allValues).Output()
        if err != nil {
            logs.Error("DeployStapService deploy network-socket getting socat error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", "Error starting tcpreplay for network->socket service")
            return err
        }
        pidValueSocat := strings.Split(string(pid), "\n")

        //check for pid error
        if pidValueSocat[0] == "" {
            pipeErrors, _ := ioutil.ReadAll(stdError)
            logs.Error(string(pipeErrors))
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", string(pipeErrors))
            return errors.New("Deploy network-socket error: Please, check warning log")
        }else{
            err = ndb.UpdatePluginValue(anode["service"],"pid",pidValueSocat[0])
            if err != nil {
                logs.Error("DeployStapService update DB pid to value Error: "+err.Error())
                _ = ndb.InsertPluginCommand(uuid, "status", "Error")
                _ = ndb.InsertPluginCommand(uuid, "output", "Error updating PID value for network->socket service.")
                return err
            }
        }

        //get tcpdump
        var grepTCPDUMP string
        for x := range allPlugins{
            if allPlugins[x]["type"] == "network-socket" && allPlugins[x]["tcpdump"] != "none"{
                grepTCPDUMP = grepTCPDUMP + "| grep -v "+allPlugins[x]["tcpdump"]
            }
        }

        time.Sleep(time.Second * 1) 
        DumpTCP := strings.Replace(tcpdumpPID, "<TCPDUMP>", grepTCPDUMP, -1)
        DumpBPF := strings.Replace(DumpTCP, "<IFACE>", allPlugins[anode["service"]]["interface"], -1)
        allDumpValues := strings.Replace(DumpBPF, "<BPF>", allPlugins[anode["service"]]["bpf"], -1)

        //get pid for tcpdump
        pid, err = exec.Command(command, param, allDumpValues).Output()
        if err != nil {
            logs.Error("DeployStapService deploy network-socket getting tcpdump pid error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", "Error getting tcpdump PID value for network->socket service.")
            return err
        }
        pidValueTcpdump := strings.Split(string(pid), "\n")

        //check for a new tcpdump PID
        if pidValueTcpdump[0] != "" {
            err = ndb.UpdatePluginValue(anode["service"],"tcpdump",pidValueTcpdump[0])
            if err != nil {
                logs.Error("DeployStapService update DB tcpdump to value Error: "+err.Error())
                _ = ndb.InsertPluginCommand(uuid, "status", "Error")
                _ = ndb.InsertPluginCommand(uuid, "output", "Error updating tcpdump PID value for network->socket service.")
                return err
            }
        }

        _ = ndb.InsertPluginCommand(uuid, "status", "Success")
        _ = ndb.InsertPluginCommand(uuid, "output", pidValueSocat[0]+" - "+pidValueTcpdump[0])
        logs.Notice("Deploy "+allPlugins[anode["service"]]["type"]+" successfully --> Description: "+allPlugins[anode["service"]]["name"]+"  --  SOCAT: "+pidValueSocat[0]+"  --  TCPDUMP: "+pidValueTcpdump[0])
    }
    
    return nil
}

func StopStapService(anode map[string]string)(err error) {
    allPlugins,err := ndb.GetPlugins()
    if err != nil {logs.Error("Error! can't read database for stop the service: "+err.Error())}

    command, err := utils.GetKeyValueString("execute", "command")
    if err != nil { logs.Error("Error getting data from main.conf")}
    param, err := utils.GetKeyValueString("execute", "param")
    if err != nil { logs.Error(" Error getting data from main.conf")}
    socatPID, err := utils.GetKeyValueString("execute", "socatPID")
    if err != nil { logs.Error(" Error getting data from main.conf")}
    openSSL, err := utils.GetKeyValueString("execute", "openSSL")
    if err != nil { logs.Error(" Error getting data from main.conf")}
    tcpdumpPID, err := utils.GetKeyValueString("execute", "tcpdumpPID")
    if err != nil { logs.Error(" Error getting data from main.conf")}
    
    //insert common values into command db
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated); if err != nil {logs.Error("StopStapService Error inserting output into database")}
    _ = ndb.InsertPluginCommand(uuid, "id", anode["service"]); if err != nil {logs.Error("StopStapService Error inserting identifier into database")}
    _ = ndb.InsertPluginCommand(uuid, "type", allPlugins[anode["service"]]["type"]); if err != nil {logs.Error("StopStapService Error inserting type into database")}
    _ = ndb.InsertPluginCommand(uuid, "action", "Stop"); if err != nil {logs.Error("StopStapService Error inserting type into database")}
    
    if allPlugins[anode["service"]]["type"] == "socket-network" {
        //get pid
        pid, err := exec.Command(command, param, strings.Replace(socatPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)).Output()
        if err != nil {
            logs.Error("StopStapService deploy socket STAP Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error"); if err != nil {logs.Error("StopStapService Error inserting status into database")}
            _ = ndb.InsertPluginCommand(uuid, "output", "Error getting PID value for socket->network service."); if err != nil {logs.Error("StopStapService Error inserting output into database")}
            return err
        }
        
        //insert command values
        _ = ndb.InsertPluginCommand(uuid, "command", strings.Replace(socatPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)); if err != nil {logs.Error("StopStapService Error inserting command into database")}
        _ = ndb.InsertPluginCommand(uuid, "description", "Stop socket->network service"); if err != nil {logs.Error("StopStapService Error inserting description into database")}
        
        pidValue := strings.Split(string(pid), "\n")
        //Killing PID
        for z := range pidValue{
            pidToInt,_ := strconv.Atoi(pidValue[z])
            process, _ := os.FindProcess(pidToInt)
            _ = process.Kill()
        }
    
        //Updating service PID
        err = ndb.UpdatePluginValue(anode["service"],"pid","none")
        if err != nil {
            logs.Error("StopStapService update DB pid to none Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error"); if err != nil {logs.Error("StopStapService Error inserting status into database")}
            _ = ndb.InsertPluginCommand(uuid, "output", "Error deleting PID value after socket->network service stops."); if err != nil {logs.Error("StopStapService Error inserting output into database")}
            return err
        }
    }else if allPlugins[anode["service"]]["type"] == "socket-pcap" {
        pid, err := exec.Command(command, param, strings.Replace(socatPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)).Output()
        if err != nil {
            logs.Error("StopStapService deploy socket STAP Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error"); if err != nil {logs.Error("StopStapService Error inserting status into database")}
            _ = ndb.InsertPluginCommand(uuid, "output", "Error getting PID value for socket->network service."); if err != nil {logs.Error("StopStapService Error inserting output into database")}
            return err
        }
        
        //insert command values
        _ = ndb.InsertPluginCommand(uuid, "command", strings.Replace(socatPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)); if err != nil {logs.Error("StopStapService Error inserting command into database")}
        _ = ndb.InsertPluginCommand(uuid, "description", "Stop socket->pcap service"); if err != nil {logs.Error("StopStapService Error inserting description into database")}
        
        pidValue := strings.Split(string(pid), "\n")
        //killing PID
        for z := range pidValue{
            pidToInt,_ := strconv.Atoi(pidValue[z])
            process, _ := os.FindProcess(pidToInt)
            _ = process.Kill()
        }
    
        //Updating service PID
        err = ndb.UpdatePluginValue(anode["service"],"pid","none")
        if err != nil {
            logs.Error("StopStapService update DB pid to none Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error"); if err != nil {logs.Error("StopStapService Error inserting status into database")}
            _ = ndb.InsertPluginCommand(uuid, "output", "Error deleting PID value after socket->pcap service stops."); if err != nil {logs.Error("StopStapService Error inserting output into database")}
            return err
        }
    }else if allPlugins[anode["service"]]["type"] == "network-socket" {
        //insert command values
        _ = ndb.InsertPluginCommand(uuid, "command", strings.Replace(socatPID, "<PORT>", allPlugins[anode["service"]]["port"], -1)); if err != nil {logs.Error("StopStapService Error inserting command into database")}
        _ = ndb.InsertPluginCommand(uuid, "description", "Stop network->socket service"); if err != nil {logs.Error("StopStapService Error inserting description into database")}

        //kill OPENSSL
        collector := strings.Replace(openSSL, "<COLLECTOR>",allPlugins[anode["service"]]["collector"], -1)
        allValues := strings.Replace(collector, "<PORT>", allPlugins[anode["service"]]["port"], -1)
        pid, err := exec.Command(command, param, allValues).Output()
        if err != nil {
            logs.Error("StopStapService deploy network-socket STAP Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error"); if err != nil {logs.Error("StopStapService Error inserting status into database")}
            _ = ndb.InsertPluginCommand(uuid, "output", "Error getting PID value for network->socket service."); if err != nil {logs.Error("StopStapService Error inserting output into database")}
            return err
        }
        pidValue := strings.Split(string(pid), "\n")
        for z := range pidValue{
            pidToInt,_ := strconv.Atoi(pidValue[z])
            process, _ := os.FindProcess(pidToInt)
            _ = process.Kill()
        }

        //Kill TCPDUMP
        var grepTCPDUMP string
        for x := range allPlugins{
            if allPlugins[x]["type"] == "network-socket" && allPlugins[x]["tcpdump"] != "none"{
                grepTCPDUMP = grepTCPDUMP + "| grep -v "+allPlugins[x]["tcpdump"]
            }
        }
        DumpTCP := strings.Replace(tcpdumpPID, "<TCPDUMP>", grepTCPDUMP, -1)
        DumpBPF := strings.Replace(DumpTCP, "<IFACE>", allPlugins[anode["service"]]["interface"], -1)
        allDumpValues := strings.Replace(DumpBPF, "<BPF>", allPlugins[anode["service"]]["bpf"], -1)
        pid, err = exec.Command(command, param, allDumpValues).Output()
        if err != nil {
            logs.Error("StopStapService deploy network-socket getting tcpdump pid error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error"); if err != nil {logs.Error("StopStapService Error inserting status into database")}
            _ = ndb.InsertPluginCommand(uuid, "output", "Error getting tcpdump PID value for kill network->socket service."); if err != nil {logs.Error("StopStapService Error inserting output into database")}
            return err
        }
        pidValueTcpdump := strings.Split(string(pid), "\n")
        for v := range pidValueTcpdump{
            pidToInt,_ := strconv.Atoi(pidValueTcpdump[v])
            process, _ := os.FindProcess(pidToInt)
            _ = process.Kill()
        }

        err = ndb.UpdatePluginValue(anode["service"],"tcpdump","none")
        if err != nil {
            logs.Error("StopStapService update DB tcpdump to none Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error"); if err != nil {logs.Error("StopStapService Error inserting status into database")}
            _ = ndb.InsertPluginCommand(uuid, "output", "Error deleting tcpdump PID value after network->socket service stops."); if err != nil {logs.Error("StopStapService Error inserting output into database")}
            return err
        }
        err = ndb.UpdatePluginValue(anode["service"],"pid","none")
        if err != nil {
            logs.Error("StopStapService update DB pid to none Error: "+err.Error())
            _ = ndb.InsertPluginCommand(uuid, "status", "Error"); if err != nil {logs.Error("StopStapService Error inserting status into database")}
            _ = ndb.InsertPluginCommand(uuid, "output", "Error deleting PID value after network->socket service ."); if err != nil {logs.Error("StopStapService Error inserting output into database")}
            return err
        }
    }
    _ = ndb.InsertPluginCommand(uuid, "status", "Success"); if err != nil {logs.Error("StopStapService Error inserting status into database")}
    logs.Notice(allPlugins[anode["service"]]["type"]+" service stopped successfuly!")

    return nil
}

func ChangeSuricataTable(anode map[string]string)(err error) {
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
    
    return nil
}

func GetServiceCommands(anode map[string]string)(data map[string]map[string]string, err error) {
    data, err = ndb.GetPluginCommands()
    commandService := make(map[string]map[string]string)
    for id,val := range data{
        if data[id]["id"] == anode["service"]{
            commandService[id] = map[string]string{}
            for val,_ := range val{
                commandService[id][val] = data[id][val]
            }
        }
    }
    if err != nil {logs.Error("GetServiceCommands Error: "+err.Error()); return nil, err}

    return commandService, err
}

func StopPluginsGracefully()(){
    command, err := utils.GetKeyValueString("execute", "command")
    if err != nil { logs.Error("Error getting data from main.conf")}
    param, err := utils.GetKeyValueString("execute", "param")
    if err != nil { logs.Error(" Error getting data from main.conf")}
    socatPID, err := utils.GetKeyValueString("execute", "socatPID")
    if err != nil { logs.Error(" Error getting data from main.conf")}
    openSSL, err := utils.GetKeyValueString("execute", "openSSL")
    if err != nil { logs.Error(" Error getting data from main.conf")}
    tcpdumpPID, err := utils.GetKeyValueString("execute", "tcpdumpPID")
    if err != nil { logs.Error(" Error getting data from main.conf")}
    suricataBackup, err := utils.GetKeyValueString("suricata", "backup")  
    if err != nil {logs.Error("StopSuricataService Error getting data from main.conf: "+err.Error())}
    suricataPidfile, err := utils.GetKeyValueString("suricata", "pidfile")  
    if err != nil {logs.Error("StopSuricataService Error getting data from main.conf: "+err.Error())}
    plugins,err := ndb.GetPlugins()
    if err != nil {logs.Error("StopPluginsGracefully Error: "+err.Error())}

    for id := range plugins{
        if plugins[id]["type"] == "suricata"{
            if plugins[id]["pid"] != "none"{
                //kill suricata process
                PidInt,_ := strconv.Atoi(strings.Trim(string(plugins[id]["pid"]), "\n"))
                process, _ := os.FindProcess(PidInt)
                _ = process.Kill()
                //delete pid file
                _ = os.Remove(suricataBackup+id+"-"+suricataPidfile)
            }            
        }else if plugins[id]["type"] == "socket-network" || plugins[id]["type"] == "socket-pcap" {
            if plugins[id]["pid"] != "none"{
                pid, _ := exec.Command(command, param, strings.Replace(socatPID, "<PORT>", plugins[id]["port"], -1)).Output()
                pidValue := strings.Split(string(pid), "\n")
                //Killing PID
                for z := range pidValue{
                    pidToInt,_ := strconv.Atoi(pidValue[z])
                    process, _ := os.FindProcess(pidToInt)
                    _ = process.Kill()
                }
            }
        }else if plugins[id]["type"] == "network-socket"{
            if plugins[id]["pid"] != "none" || plugins[id]["tcpdump"] != "none"{
                //kill tcpreplay
                collector := strings.Replace(openSSL, "<COLLECTOR>",plugins[id]["collector"], -1)
                allValues := strings.Replace(collector, "<PORT>", plugins[id]["port"], -1)
                pid, _ := exec.Command(command, param, allValues).Output()
                pidValue := strings.Split(string(pid), "\n")
                for z := range pidValue{
                    pidToInt,_ := strconv.Atoi(pidValue[z])
                    process, _ := os.FindProcess(pidToInt)
                    _ = process.Kill()
                }
                //kill tcpdump
                var grepTCPDUMP string
                for x := range plugins{
                    if plugins[x]["type"] == "network-socket" && plugins[x]["tcpdump"] != "none"{
                        grepTCPDUMP = grepTCPDUMP + "| grep -v "+plugins[x]["tcpdump"]
                    }
                }
                DumpTCP := strings.Replace(tcpdumpPID, "<TCPDUMP>", grepTCPDUMP, -1)
                DumpBPF := strings.Replace(DumpTCP, "<IFACE>", plugins[id]["interface"], -1)
                allDumpValues := strings.Replace(DumpBPF, "<BPF>", plugins[id]["bpf"], -1)
                pid, _ = exec.Command(command, param, allDumpValues).Output()
                pidValueTcpdump := strings.Split(string(pid), "\n")
                for v := range pidValueTcpdump{
                    pidToInt,_ := strconv.Atoi(pidValueTcpdump[v])
                    process, _ := os.FindProcess(pidToInt)
                    _ = process.Kill()
                }
            }
        }
    }
    //kill zeek
    err = zeek.StopingZeek()
    if err != nil {logs.Error("StopPluginsGracefully Error stopping Zeek: "+err.Error())}
}